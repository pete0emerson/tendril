#!/bin/bash

if [ "$1" == "help" ] ; then
	cat << EOF
---
short: Show my favorite vegetable
long: 'Show my favorite vegetable


This is a trivial example of a tendril command that shows my favorite vegetable
EOF
	exit 0
fi

echo "onion"
