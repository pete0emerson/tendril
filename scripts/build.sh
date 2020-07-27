#!/bin/bash

set -ex

buildTime=$(date +%s)
buildDate=$(date -d @${buildTime} -u "+%Y-%m-%d %H:%M:%S UTC")
buildHash=$(git rev-parse HEAD)
if [ ! -z "$(git status --porcelain)" ]; then
        buildHash="$buildHash (with uncommitted changes)"
fi
buildVersion="development"
buildOS=$(uname -s | tr 'A-Z' 'a-z')
buildInstallMethod=make

mkdir -p bin

for app in $(ls cmd) ; do
	cd cmd/$app
	go build\
		-ldflags "-X 'main.buildDate=${buildDate}'
		-X 'main.buildHash=${buildHash}'
		-X 'main.buildVersion=${buildVersion}'
		-X 'main.buildOS=${buildOS}'
		-X 'main.buildInstallMethod=${buildInstallMethod}'" \
		-o ../../bin/$app
	cd ..
done
