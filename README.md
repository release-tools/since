# changelog-parser

Parse changelog file and return changes for a given version.

## Installation

```bash
go install github.com/outofcoffee/changelog-parser
```

## Usage

```bash
changelog-parser [--changelog CHANGELOG.md] [--version 0.1.0]
```

## Example

```bash
$ cat CHANGELOG.md
# Change Log
...
## [0.1.0] - 2023-03-04
### Added
- feat: some change.
- feat: another change.

$ changelog-parser --version 0.1.0
## [0.1.0] - 2023-03-04
### Added
- feat: some change.
- feat: another change.
```

If no version is specified, the latest version is returned.
