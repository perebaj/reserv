GOLANGCI_LINT_VERSION=v1.62.2
GOLANG_VULCHECK_VERSION=v1.1.3

## run all tests. Usage `make test` or `make test testcase="TestFunctionName"` to run an isolated tests
.PHONY: test
test:
	if [ -n "$(testcase)" ]; then \
		go test ./... -timeout 10s -race -run="^$(testcase)$$" -v; \
	else \
		go test ./... -timeout 10s -race; \
	fi

## Run nil-checker to find potential nil pointer dereferences
PHONY: nil-checker
nil-checker:
	go run go.uber.org/nilaway/cmd/nilaway -test=false -include-pkgs="github.com/perebaj/reserv" -exclude-errors-in-files=mock_ ./...

## Run linter
.PHONY: lint
lint: nil-checker
	@echo "Running linter..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./... -v --timeout 5m
	go run golang.org/x/vuln/cmd/govulncheck@$(GOLANG_VULCHECK_VERSION) ./...

## Run test coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Display help for all targets
.PHONY: help
help:
	@awk '/^.PHONY: / { \
		msg = match(lastLine, /^## /); \
			if (msg) { \
				cmd = substr($$0, 9, 100); \
				msg = substr(lastLine, 4, 1000); \
				printf "  ${GREEN}%-30s${RESET} %s\n", cmd, msg; \
			} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
