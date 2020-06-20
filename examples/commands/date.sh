#!/bin/bash

if [ "$1" == "help" ] ; then
	cat << EOF
---
short: Run the Date command
long: 'Run the date command


This is a trivial example of a baton command

This will simply run the unix \`date\` command.'
EOF
	exit 0
fi

date $1
