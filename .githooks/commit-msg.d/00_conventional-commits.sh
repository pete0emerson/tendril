#!/bin/bash

echo -n "### Checking conventional commits ... "

MESSAGE=$(head -n1 $1)
regex='^(build|ci|docs|feat|fix|perf|refactor|style|test|chore|release)(\([0-9a-z-]+\))?: '

if [[ ! $MESSAGE =~ $regex ]] ; then
	cat <<EOF
failed

Commit message does not follow Conventional Commits format (https://www.conventionalcommits.org):

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]

This regular expression must pass: $regex

########################################

$MESSAGE
EOF
	exit 1
fi
echo "okay"
exit 0
