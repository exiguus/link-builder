#!/bin/bash

# Run "make qlty" before push
if ! make qlty; then
	echo "Pre-psh hook: 'make qlty' failed. Commit aborted."
	exit 1
fi

# Run "make coverage" before committing
if ! make coverage; then
	echo "Pre-commit hook: 'make coverage' failed. Commit aborted."
	exit 1
fi

# Allow the push to proceed
exit 0
