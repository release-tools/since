#!/usr/bin/env bash
set -e

NEXT_VERSION="$( since project version -q )"

NEW_CHANGELOG="$( mktemp -t changelog.XXXXXX )" || exit 1
since changelog update -q > "$NEW_CHANGELOG"
mv "$NEW_CHANGELOG" CHANGELOG.md

git add CHANGELOG.md
git commit -m "build: release ${NEXT_VERSION}."
git tag "${NEXT_VERSION}"

echo "Release ${NEXT_VERSION} created."
