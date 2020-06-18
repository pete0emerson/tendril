build:	## Build the binaries
	@scripts/build.sh

help:	## Show this help

.PHONY: build help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

#release:	## Create VERSION, CHANGELOG, tag, and release
#	@scripts/release.sh
