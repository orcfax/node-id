.DEFAULT_GOAL := help

.PHONY: lint snapshot

lint:                                     ## Lint the source code
	staticcheck ./...
	go fmt ./...
	golint ./...
	go vet ./...

binary-skip-release:					  ## Create a binary but skip the release
	goreleaser release --skip-publish --clean

binary-release: 						  ## Create a binary and perform the release
	goreleaser release --skip-publish --clean

snapshot:					  			  ## Create a snapshot binary
	goreleaser release --snapshot --skip-publish --clean

pre-commit-checks:                        ## Run pre-commit-checks
	pre-commit run --all-files

verify-signatures:                        ## Verify all signing
	./verify_sig.sh

verify-checksum-signature:                ## Verify checksum signing (useful for outputting GPG info)
	gpg --verify dist/checksums_sha256.txt.sig

help:                                                       ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
