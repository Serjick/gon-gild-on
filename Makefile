VERBOSE?=@
BIN_DIR?=${HOME}/.local/bin
GOLANGCI_LINT_VERSION?=2.0.2# NOTE: on every version bump don't forget to update ./.golangci.yml with configs from https://github.com/golangci/golangci-lint/blob/v${LINTER_VERSION}/.golangci.reference.yml

default:

test:
	${VERBOSE} go test --race --vet= --failfast --count=1 --timeout=5m --test.run= --covermode=atomic --coverprofile=coverage.out --coverpkg=./... ./... -v
	${VERBOSE} go test --fuzz ^Fuzz --fuzztime=10s ./golden/gildedsergigodiff -v

coverage: test
	${VERBOSE} go tool cover --func=coverage.out

lint: golangci-lint
	${VERBOSE} ${GOLANGCI_LINT} run --timeout=2m ./... -v

.PHONY: test coverage lint

golangci-lint:
ifneq (, $(wildcard ${BIN_DIR}/golangci-lint@${GOLANGCI_LINT_VERSION}/golangci-lint))
GOLANGCI_LINT=${BIN_DIR}/golangci-lint@${GOLANGCI_LINT_VERSION}/golangci-lint
else
ifneq (, $(shell which golangci-lint))
GOLANGCI_LINT=$(shell which golangci-lint)
else
	${VERBOSE} GOBIN=${BIN_DIR}/golangci-lint@${GOLANGCI_LINT_VERSION} go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}
GOLANGCI_LINT=${BIN_DIR}/golangci-lint@${GOLANGCI_LINT_VERSION}/golangci-lint
endif
endif

guard-%: GUARD
	@if [ -z '${${*}}' ]; then echo 'Variable $* not set.' && exit 1; fi

GUARD:

.PHONY: GUARD

%:
	@:
