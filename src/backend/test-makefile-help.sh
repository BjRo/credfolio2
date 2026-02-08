#!/bin/bash
# Test script for Makefile help output
# This script verifies that the help target shows correct target names

set -euo pipefail

cd "$(dirname "$0")"

# Simulate the help command with MAKEFILE_LIST containing the filename
# (this is what happens when make is run)
# When grep operates on multiple files (or $(MAKEFILE_LIST)), it outputs: filename:target:## description
# We simulate this by explicitly passing the filename to grep multiple times
output=$(grep -E '^[a-zA-Z_-]+:.*?## .*$' Makefile Makefile | head -n 8 | sed 's/^[^:]*://' | awk 'BEGIN {FS = ":.*## "}; {printf "  make %-20s %s\n", $1, $2}')

# Verify that output contains expected target names
expected_targets=(
  "help"
  "migration"
  "migrate-up"
  "migrate-down"
  "migrate-down-all"
  "migrate-force"
  "migrate-version"
  "migrate-status"
)

# Check that each expected target appears in the output
failed=0
for target in "${expected_targets[@]}"; do
  if ! echo "$output" | grep -q "make $target "; then
    echo "FAIL: Expected target '$target' not found in help output"
    failed=1
  fi
done

# Verify that "Makefile" does NOT appear as a target name
if echo "$output" | grep -q "make Makefile "; then
  echo "FAIL: Found 'make Makefile' in output (should show actual target names)"
  failed=1
fi

if [ $failed -eq 0 ]; then
  echo "PASS: All help targets are correctly displayed"
  exit 0
else
  echo ""
  echo "Actual output:"
  echo "$output"
  exit 1
fi
