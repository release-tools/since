---
name: since
description: >-
  Use the `since` tool to manage releases, changelogs and versioning driven by
  git history and conventional commits. Use when asked to perform a release, cut
  a version, generate or update a CHANGELOG, work out the next semantic version,
  list changes since the last release, or scaffold a since.yaml config. Keywords:
  since, release, changelog, semver, next version, conventional commits, tag.
---

# since

`since` parses git history and conventional commits to drive releases. It works
out the next [semantic version](http://semver.org), generates changelog entries
from the commits since the last tag, and can commit and tag a release in one
step.

Commands are grouped under `changelog`, `project`, and a top-level `init`. Run
`since <command> --help` for the full flag list. Install it with
`brew install release-tools/tap/since` or `go install github.com/release-tools/since`.

## Perform a release

This is the headline task. "Perform a release" (also: "cut a release", "release
the project", "tag a new version") means run:

```bash
since project release
```

This generates a new changelog section from the commits since the last release,
**commits** the updated `CHANGELOG.md`, and creates a **new git tag** for the
computed version — all in one step. It does not push; push the branch and tags
separately once you've confirmed the result.

Before running it, check for a `since.yaml` config (see below) — it may enforce
a required branch and run `before`/`after` hooks (tests, publish steps). A
failing hook aborts the release.

Useful flags:

- `-c, --changelog` — path to the changelog file (default `CHANGELOG.md`).
- `-g, --git-repo` — path to the repository (default `.`).
- `-o, --order-by` — how the latest tag is chosen: `semver` (default),
  `commit-date`, or `alphabetical`.
- `-t, --tag` — include commits after this specific tag.
- `--unique` — de-duplicate commit messages (default true).

Typical flow:

```bash
since project version        # preview the version that will be cut
since project changes        # review what will go into the changelog
since project release        # commit the changelog and create the tag
git push --follow-tags       # publish the commit and tag
```

## Scaffold a config file: `since init`

```bash
since init
```

Creates a `since.yaml` in the current directory, pre-populated with commented
examples for branch requirements, pre/post-release hook scripts, and commit
exclusions. If the file already exists it is **overwritten**. Write it elsewhere
with `-o, --output <dir>`.

Uncomment the parts you need. Key settings:

- `requireBranch` — refuse to run unless on this branch (exact or glob, e.g.
  `release/*`). Checked automatically before any command.
- `ignore` — patterns matched against commit subject lines; matching commits are
  left out of the changelog (e.g. `chore:`, `docs:`, `Merge pull request`).
- `before` / `after` — hooks run in order around the release; any failure aborts
  it. Define each as either `command` + `args`, or an inline `script`.

Hooks receive these environment variables: `SINCE_NEW_VERSION`,
`SINCE_OLD_VERSION`, `SINCE_SHA`, `SINCE_REPO_PATH`.

## Inspect changes and versions (no writes)

- `since project changes` — list commits since the last tag, grouped by
  conventional-commit type.
- `since project version` — print the next semantic version implied by those
  commits. Add `-c, --current` to print the current version instead.

Both accept `-g/--git-repo`, `-o/--order-by`, and `-t/--tag`.

## Work with changelog files directly

- `since changelog generate` — build a new changelog from an existing one plus
  the latest commits, and print it to stdout (doesn't write in place).
- `since changelog update` — write the new release section into the existing
  changelog file.
- `since changelog extract` — pull out the entries for one version (`-v`, or the
  most recent if omitted). `--header` includes the version heading. Handy for
  feeding release notes to a GitHub Release.
- `since changelog init` — create a fresh changelog file from the repo's git
  history.

Global flags on these: `-c/--changelog` (file path), `--output-file` (write
somewhere other than stdout), `-l/--log-level`, and `-q/--quiet` (silence logs
for scripting).

## Notes

- Version bumps follow conventional commits: `feat:` → minor; `fix:` and other
  types (`build`, `chore`, `ci`, `docs`, `refactor`, `security`, `style`,
  `test`) → patch; a `!` suffix (e.g. `feat!:`) or a `BREAKING CHANGE` footer →
  major.
- Use `-q/--quiet` when capturing output in scripts so log lines don't pollute
  stdout.
- In CI, the `release-tools/since` GitHub Action can extract the latest release
  notes into a file for `softprops/action-gh-release`.
