# Example: Print changes

Given a changelog file:

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
```

Print the changes for version `0.1.0`:

```bash
$ since changelog extract --version 0.1.0
## [0.1.0] - 2023-03-04
### Added
- feat: another change.
```

If no version is specified, the latest version is returned.
