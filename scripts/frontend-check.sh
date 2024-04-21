#!/bin/bash

# navigate to frontend directory
cd frontend

echo "Running npm format..."
npm run format
FORMAT_STATUS=$?

echo "Running npm lint..."
npm run lint
LINT_STATUS=$?

echo "Running npm check..."
npm run check
CHECK_STATUS=$?

# navigate back to root directory
cd ..

# check if any of the npm command exited with a non-zero status
if [ $FORMAT_STATUS -ne 0 -o $LINT_STATUS -ne 0 -o $CHECK_STATUS -ne 0 ]; then
    echo "One or more checks failed. Please fix the errors before committing."
    exit 1
fi

echo "All checks passed. Ready to commit."
exit 0
