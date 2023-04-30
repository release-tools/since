# since
[![CI](https://github.com/outofcoffee/since/actions/workflows/ci.yaml/badge.svg)](https://github.com/outofcoffee/since/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/outofcoffee/since)](https://goreportcard.com/report/github.com/outofcoffee/since)
[![Go Reference](https://pkg.go.dev/badge/github.com/outofcoffee/since.svg)](https://pkg.go.dev/github.com/outofcoffee/since)
![License](https://img.shields.io/github/license/outofcoffee/since)

- Parse git history and generate changelog.
- Calculate the next version based on [semver](http://semver.org) and [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).
- Parse changelog files and extract changes for a given version.

## Installation

### Homebrew

```bash
brew install outofcoffee/tap/since
```

### Go

```bash
go install github.com/outofcoffee/since
```

## Usage

### `changelog update`

Generates a new changelog based on an existing changelog file,
using the commits since the last release.

```bash
Usage:
  since changelog update [flags]

Flags:
  -c, --changelog string   Path to changelog file (default "CHANGELOG.md")
  -g, --git-repo string    Path to git repository (default ".")
  -h, --help               help for update
  -o, --order-by string    How to determine the latest tag (alphabetical|commit-date|semver)) (default "semver")

Global Flags:
  -l, --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "debug")
  -q, --quiet              Whether to disable logging
```

---

### `changelog extract`

Extracts changes for a given version in a changelog file.
If no version is specified, the most recent version is used.

```bash
Usage:
  since changelog extract [flags]

Flags:
      --header           whether to include the version header in the output
  -h, --help             help for extract
  -v, --version string   Version to parse changelog for

Global Flags:
  -c, --changelog string   Path to changelog file (default "CHANGELOG.md")
  -l, --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "debug")
  -q, --quiet              Whether to disable logging
```

---

### `project changes`

Reads the commit history for the current git repository, startingnce project changes -h
from the most recent tag. Lists the commits categorised by their type.

```bash
Usage:
  since project changes [flags]

Flags:
  -h, --help   help for changes

Global Flags:
  -g, --git-repo string    Path to git repository (default ".")
  -l, --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "debug")
  -o, --order-by string    How to determine the latest tag (alphabetical|commit-date|semver)) (default "semver")
  -q, --quiet              Whether to disable logging
  -t, --tag string         Include commits after this tag
```

---

### `project version`

Reads the commit history for the current git repository, startingnce project version -h
from the most recent tag. Returns the next semantic version
based on the changes.

Changes influence the version according to
[conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)

```bash
Usage:
  since project version [flags]

Flags:
  -c, --current   Just print the current version
  -h, --help      help for version

Global Flags:
  -g, --git-repo string    Path to git repository (default ".")
  -l, --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "debug")
  -o, --order-by string    How to determine the latest tag (alphabetical|commit-date|semver)) (default "semver")
  -q, --quiet              Whether to disable logging
  -t, --tag string         Include commits after this tag
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

$ since changelog extract --version 0.1.0
## [0.1.0] - 2023-03-04
### Added
- feat: another change.
```

If no version is specified, the latest version is returned.
