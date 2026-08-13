#!/bin/sh
set -e
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
git config core.hooksPath .githooks
echo "Git hooks enabled from: $PROJECT_ROOT/.githooks"
echo "Active hooks: commit-msg (Conventional Commits validation), pre-push (go vet + npm lint when present in repo)"
