# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.7.0] - 2023-04-30
### Added
- feat: allows printing of current version.

### Changed
- build: adds release script.
- ci: renames tap repo.

### Fixed
- fix: removes redundant whitespace.
- fix: prints directly to stdout.

## [0.6.1] - 2023-04-30
### Changed
- build: ignores binary name.
- ci: updates dryrun changelog command.
- ci: replaces deprecated goreleaser flag.

## [0.6.0] - 2023-04-30
### Added
- feat: groups changes into sections.
- feat: adds changelog update command.
- feat: adds command to list changes in repo since tag.

### Changed
- ci: adds goreleaser config and workflow step.
- refactor: organises commands under 'project' and 'changelog'.
- refactor: renames list command to extract.
- build: renames module.

### Fixed
- fix: only fetch commits once when updating changelog.
- fix: sorts commit categories before printing.

## [0.5.0] - 2023-04-28
### Added
- feat: determine semver based on git log.

## [0.4.1] - 2023-04-28
### Changed
- refactor: switches commander to cobra.

## [0.4.0] - 2023-03-05
### Added
- feat: allow version header to be skipped.

## [0.3.0] - 2023-03-05
### Changed
- refactor: moves printer to function parameter.

## [0.2.0] - 2023-03-05
### Added
- fix: prints entries to stdout.
- ci: adds GitHub Actions workflow.

## [0.1.0] - 2023-03-04
### Added
- feat: initial version.
