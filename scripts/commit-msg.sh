#!/bin/bash
REGEX="^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9-]+\))?!?: .+"

FILE=`cat $1` # File containing the commit message

echo "Commit Message: ${FILE}"

if ! [[ $FILE =~ $REGEX ]]; then
	echo "Invalid commit message."
	echo "  Please follow the Conventional Commit standard."
	echo "  Commit messages start with: (build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)."
	echo ""
	echo "  Example: feat: add new feature"
	echo "  Example: fix: fix a bug"
	echo "  Example: ci(build): adjust build process"
	echo ""
	echo "  See: <https://www.conventionalcommits.org/en/v1.0.0/#specification>."
	exit 1
else
	echo "Valid commit message."
fi