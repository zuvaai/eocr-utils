# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

### Added

- `pkg/eocr` now has a `NewDocumentFromText` function to create documents from a UTF8 text.

## [v0.0.1]

### Added

- `recognition_results.proto` and generated type and protobuf code are in the package `pkg/document`
- `pkg/eocr` contains the main marshall/unmarshall code for the eocr file format
- `pkg/eocr` contains the all the errors (variables begining with `Err`)
- CODEOWNERS configuration
- README
- Release command `cmd/release`

### Removed

- Headers and Footers have been removed from `recognition_results.proto` (labels 9 and 10)

### Changed

- Some versions of `recognition_results.proto` used the name `CharacterRange` for spans. In this repo, we use the name `Span`. Since protobuf does structural matching, this does not impact interoperability.
