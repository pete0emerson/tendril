#!/bin/bash

if [ "$1" == "help" ] ; then
	cat << EOF
---
short: Show my favorite fruit
long: 'Show my favorite fruit


This is a trivial example of a tendril command that shows my favorite fruit
EOF
	exit 0
fi

echo "kiwi"
