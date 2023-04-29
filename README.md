# since

- Parses git log and generates changelog entries.
- Calculates the next version based on semver and conventional commits.
- Parses changelog files and extract changes for a given version.

## Installation

```bash
go install github.com/outofcoffee/since
```

## Usage

```bash
since extract [--changelog CHANGELOG.md] [--version 0.1.0]
```

## Example

```bash
$ cat CHANGELOG.md
# Change Log
...
## [0.2.0] - 2023-03-05
### Added
- feat: some change.

## [0.1.0] - 2023-03-04
### Added
- feat: another change.

$ since extract --version 0.1.0
## [0.1.0] - 2023-03-04
### Added
- feat: another change.
```

If no version is specified, the latest version is returned.
