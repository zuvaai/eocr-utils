package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

var Main = &cobra.Command{Use: "release", SilenceUsage: true}

func init() {
	Main.AddCommand(
		CreateCommand(),
	)
}

func CreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create <release number>",
		Short: "Create a release",
		Long:  "Create a release",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			releaseNumber := args[0]
			return createRelease(releaseNumber)
		},
	}
}

// nolint: gocyclo
func createRelease(tag string) error {
	// Check CHANGES.md
	fmt.Println("Checking CHANGES.md")
	changes, err := os.Open("./CHANGES.md")
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("release command must be run from the root directory containig a CHANGES.md file")
	} else if err != nil {
		return fmt.Errorf("release command failed: %w", err)
	}
	defer changes.Close()
	var version *semver.Version
	scanner := bufio.NewScanner(changes)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "##") {
			continue
		}
		changesVersion := strings.ToLower(strings.Trim(strings.TrimPrefix(line, "##"), " \n\r[]()"))
		if strings.Contains(changesVersion, "unreleased") {
			return fmt.Errorf("version set to 'unreleased' in CHANGES.md")
		} else if !strings.HasPrefix(changesVersion, "v") {
			return fmt.Errorf("version in CHANGES.md (%s) does not begin with v", changesVersion)
		}

		version, err = semver.NewVersion(changesVersion)
		if err != nil {
			return fmt.Errorf("version (%s) in CHANGES.md is invalid", changesVersion)
		}
		break
	}

	if version == nil {
		return fmt.Errorf("could not find version in CHANGES.md")
	}

	// Check tag
	if !strings.HasPrefix(tag, "v") {
		return fmt.Errorf("release tag must start with 'v'")
	}
	if version.Original() != tag {
		return fmt.Errorf("tag and CHANGES.md version do not match: %s != %s", tag, version)
	}
	// Execute the git command and verify that we are currently on the main branch
	output, err := runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}
	if output != "main" {
		return fmt.Errorf("must be on main branch to create release")
	}
	// check that the branch is clean
	output, err = runGitCommand("status", "--porcelain")
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("branch must be clean to create release")
	}

	// check that branch is not behind or ahead of origin/main
	output, err = runGitCommand("status", "--porcelain", "--branch")
	if err != nil {
		return err
	}
	if strings.Contains(output, "behind") {
		return fmt.Errorf("branch is behind origin/main")
	}

	if strings.Contains(output, "ahead") {
		return fmt.Errorf("branch is ahead of origin/main")
	}

	// Check that there is no tag with the same name
	output, err = runGitCommand("tag", "-l", tag)
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("release tag already exists: %s", tag)
	}

	// Check that the testdata directory exists
	_, err = os.Stat("testdata")
	if err != nil {
		return fmt.Errorf("must be at the top-level directory of the repository to create release")
	}

	// Remove the testdata directory to save space.
	// We remove the testdata directory to avoid limits in
	// https://go.dev/ref/mod#zip-path-size-constraints
	output, err = runGitCommand("rm", "-rf", "testdata")
	if err != nil {
		return fmt.Errorf("unable to remove testdata folder to reduce module size: %w", err)
	}
	fmt.Println(output)

	// Commit the result
	output, err = runGitCommand("commit", "-m", "\"Release "+tag+"\"")
	if err != nil {
		return fmt.Errorf("unable to commit : %w", err)
	}
	fmt.Println(output)

	// Tag the result
	output, err = runGitCommand("tag", tag)
	if err != nil {
		return fmt.Errorf("unable create git tag for release: %w", err)
	}
	fmt.Println(output)

	// Push the result
	output, err = runGitCommand("push", "origin", tag)
	if err != nil {
		return fmt.Errorf("unable to push release tag to repository: %w", err)
	}
	fmt.Println(output)

	// Reset the branch to origin/main
	output, err = runGitCommand("reset", "--hard", "origin/main")
	if err != nil {
		return err
	}
	fmt.Println(output)

	return nil
}

func runGitCommand(args ...string) (string, error) {
	command := fmt.Sprintf("git %s", strings.Join(args, " "))
	fmt.Println("Running command:", command)
	cmd := exec.Command("/bin/sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:\n", string(output))
		return "", err
	}
	return strings.Trim(string(output), " \n"), nil
}

func main() {
	if err := Main.Execute(); err != nil {
		os.Exit(-1)
	}
}
