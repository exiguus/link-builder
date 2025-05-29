#!/bin/bash

# Run "make all" before committing
if ! make all; then
	echo "Pre-commit hook: 'make all' failed. Commit aborted."
	exit 1
fi

# Allow the commit to proceed
exit 0
