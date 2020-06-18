#!/bin/bash

echo -n "### Checking go fmt ... "

STAGED_GO_FILES=$(git diff --cached --name-only -- '*.go')

if [[ $STAGED_GO_FILES != "" ]]; then
	for file in $STAGED_GO_FILES; do
		go fmt $file
		git add $file
	done
fi

echo "okay"
