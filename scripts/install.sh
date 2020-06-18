#!/bin/bash

set -ex

apps=$(ls cmd)
for app in $apps ; do
	SOURCE=bin/$app
        if [ -z "$PREFIX" ] ; then
                if [ -d ~/bin ] ; then
                        cp $SOURCE ~/bin/$app
                elif [ -d /usr/local/bin ] ; then
                        cp $SOURCE /usr/local/bin/$app
                else
                        echo ~/bin and /usr/local/bin do not exist. Aborting.
                        exit 1
                fi
        else
                cp $SOURCE $PREFIX/$app
        fi
done
