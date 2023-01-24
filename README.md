# Utilities for eOCR files

This repo is the canonical home for utilities related to creating and
validating eOCR files.

The repo is also the canonical home for the `recognition_results.proto` file
which is the main building block of the eOCR file format.

# Building

```
make lfs-checkout
make
```

# Running tests
```
make lfs-checkout
make test
```

# Developing

- Changes since the last version must be documented in `CHANGES.md`, see https://keepachangelog.com/en/1.0.0/

# Releasing a new version of eocr-utils

- Update the release tag in `CHANGES.md`
- Build the release executable if necessary (`make -C cmd/release`)
- Run the release command with a new tag: `cmd/release/release create v0.0.1`
