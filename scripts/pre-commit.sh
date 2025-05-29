#!/bin/bash

# Run "make format" before committing
if ! make format; then
	echo "Pre-commit hook: 'make format' failed. Commit aborted."
	exit 1
fi

# Run "make test" before committing
if ! make test; then
	echo "Pre-commit hook: 'make test' failed. Commit aborted."
	exit 1
fi

# Run "make lint" before committing
if ! make lint; then
	echo "Pre-commit hook: 'make lint' failed. Commit aborted."
	exit 1
fi

# Allow the commit to proceed
exit 0
