build:	## Build the binaries
	@scripts/build.sh

install:	## Install the binaries
	@scripts/install.sh

fmt:		## Run go fmt
	go fmt cmd/tendril/*.go

help:	## Show this help

.PHONY: build help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
