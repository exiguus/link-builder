#!/bin/bash

# Run "make all" before committing
make all

# Check if "make all" was successful
if [ $? -ne 0 ]; then
  echo "Pre-commit hook: 'make all' failed. Commit aborted."
  exit 1
fi

# Allow the commit to proceed
exit 0