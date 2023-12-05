.DEFAULT_GOAL := help

.PHONY: lint snapshot help

lint:                                   ## Lint the source code
	staticcheck ./...
	go fmt ./...
	golint ./...
	go vet ./...

build-local-snapshot:                   ## Create a build for the local architecture only.
	goreleaser build  --single-target --clean --snapshot

release-test-all:					    ## Perform a release test for all architectures.
	goreleaser release --skip-publish --clean

release-test-all-snapshot:				## Create a release snapshot e.g. when there isn't a new tag.
	goreleaser release --snapshot --skip-publish --clean

pre-commit-checks:                      ## Run pre-commit-checks.
	pre-commit run --all-files

verify-signatures:                      ## Verify all signatures.
	./verify_sig.sh

verify-checksum-signature:              ## Verify checksum signing (useful for outputting GPG info).
	gpg --verify dist/checksums_sha256.txt.sig

help:                                   ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
