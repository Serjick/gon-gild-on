VERBOSE?=@
TOOLS_DIR?=./tools
default:

test:
	${VERBOSE} go test --race --vet= --failfast --count=1 --timeout=5m --test.run= --covermode=atomic --coverprofile=coverage.out --coverpkg=./... ./... -v
	${VERBOSE} go test --fuzz ^Fuzz --fuzztime=10s ./golden/gildedsergigodiff -v

coverage: test
	${VERBOSE} go tool cover --func=coverage.out

lint:
	${VERBOSE} go tool -modfile=${TOOLS_DIR}/go.mod golangci-lint config verify
	${VERBOSE} go tool -modfile=${TOOLS_DIR}/go.mod golangci-lint run --timeout=2m ./... -v

.PHONY: test coverage lint

go-get-golangci-lint:
	${VERBOSE} echo "Updating from $$(go list -modfile=${TOOLS_DIR}/go.mod -m -f "{{ .Version }}" github.com/golangci/golangci-lint/v2) ..."
	${VERBOSE} go get -modfile=${TOOLS_DIR}/go.mod -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	${VERBOSE} go list -modfile=${TOOLS_DIR}/go.mod -m -f "{{ .Version }}" github.com/golangci/golangci-lint/v2 >.golangci-lint-version
	${VERBOSE} echo "... upto $$(cat .golangci-lint-version), don't forget to keep .golangci.yml in sync with https://github.com/golangci/golangci-lint/blob/$$(cat .golangci-lint-version)/.golangci.reference.yml"

guard-%: GUARD
	@if [ -z '${${*}}' ]; then echo 'Variable $* not set.' && exit 1; fi

GUARD:

.PHONY: GUARD

%:
	@:
