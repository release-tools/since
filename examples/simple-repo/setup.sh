#!/usr/bin/env bash
set -e

REPO_DIR="$( mktemp -d -t gittest )"
cd "${REPO_DIR}"
git init

touch README.txt
git add README.txt
git commit -a -m"First commit"
git tag 0.0.1

echo "first update" > README.txt
git commit -a -m"first update"
git tag 0.1.0

echo "second update" >> README.txt
git commit -a -m"second update"
git tag 0.1.1

echo "third update" >> README.txt
git commit -a -m"third update"
git tag 0.1.2

echo "Repo created at ${REPO_DIR}"
