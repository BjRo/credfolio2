#!/bin/bash
# Extracts bean ID from current branch name
# Branch format: <type>/<bean-id>-<description>
# Example: feat/credfolio2-abc1-docker-compose -> credfolio2-abc1

set -e

BRANCH=$(git branch --show-current)

if [[ "$BRANCH" == "main" ]] || [[ "$BRANCH" == "master" ]]; then
    echo "Error: On main/master branch, no bean ID to extract" >&2
    exit 1
fi

# Match pattern: type/bean-id-description
# Bean ID format: credfolio2-XXXX (project name + nanoid)
if [[ "$BRANCH" =~ ^[a-z]+/(credfolio2-[a-zA-Z0-9]+)-.* ]]; then
    echo "${BASH_REMATCH[1]}"
elif [[ "$BRANCH" =~ ^[a-z]+/(beans-[a-zA-Z0-9]+)-.* ]]; then
    echo "${BASH_REMATCH[1]}"
else
    echo "Error: Could not extract bean ID from branch: $BRANCH" >&2
    echo "Expected format: <type>/<bean-id>-<description>" >&2
    exit 1
fi
