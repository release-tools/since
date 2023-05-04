# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.11.1] - 2023-05-04
### Changed
- build: adds since config.
- refactor: improves logging of VCS operations.

### Fixed
- fix: inherit environment when invoking hooks.

## [0.11.0] - 2023-05-04
### Added
- feat: adds release check for required branch.

### Changed
- build: release v0.11.0.
- refactor: moves commit logic to separate file.

### Fixed
- fix: corrects YAML config deserialisation tags.
- fix: makes changelog path relative to repo root when adding to index.

## [0.11.0] - 2023-05-04
### Added
- feat: adds release check for required branch.

### Changed
- refactor: moves commit logic to separate file.

### Fixed
- fix: makes changelog path relative to repo root when adding to index.

## [0.10.0] - 2023-05-04
### Added
- feat: adds support for pre- and post- release hooks.

## [0.9.0] - 2023-05-01
### Added
- feat: adds changelog generate command.
- feat: improves logging.

### Changed
- docs: updates license header.

## [0.8.2] - 2023-04-30
### Changed
- docs: improves description of project release command.

## [0.8.1] - 2023-04-30
### Changed
- build: removes redundant release script.
- docs: describes project release command.
- docs: improves quiet flag description.

## [0.8.0] - 2023-04-30
### Added
- feat: adds project release command.

### Changed
- build: sets working directory to root in release script.

## [0.7.3] - 2023-04-30
### Fixed
- fix: trim commit messages to first line only.

## [0.7.2] - 2023-04-30
### Fixed
- fix: corrects typo in command description.

## [0.7.1] - 2023-04-30
### Changed
- docs: improves readme.

### Fixed
- fix: sorts items in each category.

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
