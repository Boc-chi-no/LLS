#!/bin/bash

LATEST_TAG=$(git describe --tags --abbrev=0)
PREVIOUS_TAG=$(git describe --tags --abbrev=0 `git rev-list --tags --skip=1 --max-count=1`)

if [ -z "$PREVIOUS_TAG" ]; then
  exit 1
fi

# Extract and categorize commits
echo "## Changelog"
echo ""
echo "### 🚀 Additions"
git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%b" | grep "\[A\]" | sed 's/\[A\] */- /' | grep "." || echo "- No additions in this release"
echo ""
echo "### 🔧 Fixes"
git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%b" | grep "\[F\]" | sed 's/\[F\] */- /' | grep "." || echo "- No fixes in this release"
echo ""
echo "### 🔄 Updates"
git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%b" | grep "\[U\]" | sed 's/\[U\] */- /' | grep "." || echo "- No updates in this release"
echo ""
echo "### 🔨 Refactors"
git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%b" | grep "\[R\]" | sed 's/\[R\] */- /' | grep "." || echo "- No refactors in this release"
echo ""
echo "### 📝 Others"
(git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%b" | grep -v "\[\(A\|F\|U\|R\)\]" | grep -v "Merge " | grep -v "Squashed commit" | sed 's/^/- /' | grep "."
git log $PREVIOUS_TAG..$LATEST_TAG --pretty=format:"%s (commit %h)" | grep -v "Version Commit" | grep -v "Merge " | grep -v "Squashed commit" | sed 's/^/- /' | grep ".") || echo "- No other changes in this release"
echo ""
