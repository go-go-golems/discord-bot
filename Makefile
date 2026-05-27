.PHONY: all test build lint lintmax golangci-lint-install gosec govulncheck \
        glazed-lint-build glazed-lint goreleaser tag-major tag-minor tag-patch release bump-glazed install

all: test build

VERSION=v0.0.1
GLAZED_LINT_BIN ?= /tmp/glazed-lint
GLAZED_LINT_PKG ?= github.com/go-go-golems/glazed/cmd/tools/glazed-lint
GLAZED_VERSION ?= $(shell GOWORK=off go list -m -f '{{.Version}}' github.com/go-go-golems/glazed 2>/dev/null)
GLAZED_LINT_TOOL_VERSION ?= v1.3.5
GLAZED_LINT_FLAGS ?= -glazedclilint.allow-paths=pkg/analysis/,pkg/cli/,pkg/cmds/fields/,pkg/cmds/logging/,pkg/cmds/sources/,pkg/help/
GLAZED_LINT_DIRS ?= ./cmd/... ./internal/... ./pkg/...
GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target
GOLANGCI_LINT_VERSION ?= $(shell cat .golangci-lint-version)
GOLANGCI_LINT_BIN ?= $(CURDIR)/.bin/golangci-lint
GOLANGCI_LINT_ARGS ?= --timeout=5m ./cmd/... ./pkg/... ./internal/...
GO_GO_GOLEMS_MODULES ?= $(shell go list -m -f '{{if not .Main}}{{.Path}}{{end}}' all 2>/dev/null | grep '^github.com/go-go-golems/' | sort -u)

golangci-lint-install:
	mkdir -p $(dir $(GOLANGCI_LINT_BIN))
	GOBIN=$(dir $(GOLANGCI_LINT_BIN)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
glazed-lint-build:
	@echo "Building glazed-lint from Glazed module..."
	@if [ -n "$(GLAZED_VERSION)" ] && [ "$(GLAZED_VERSION)" != "(devel)" ]; then \
		echo "Installing $(GLAZED_LINT_PKG)@$(GLAZED_VERSION)"; \
		GOBIN=$(dir $(GLAZED_LINT_BIN)) GOWORK=off go install $(GLAZED_LINT_PKG)@$(GLAZED_VERSION) || \
			GOBIN=$(dir $(GLAZED_LINT_BIN)) GOWORK=off go install $(GLAZED_LINT_PKG)@$(GLAZED_LINT_TOOL_VERSION); \
	else \
		echo "Installing $(GLAZED_LINT_PKG)@$(GLAZED_LINT_TOOL_VERSION)"; \
		GOBIN=$(dir $(GLAZED_LINT_BIN)) GOWORK=off go install $(GLAZED_LINT_PKG)@$(GLAZED_LINT_TOOL_VERSION); \
	fi

glazed-lint: glazed-lint-build
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)


lint: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) config verify
	$(GOLANGCI_LINT_BIN) run -v $(GOLANGCI_LINT_ARGS)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

lintmax: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) run -v --max-same-issues=100 $(GOLANGCI_LINT_ARGS)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude-generated -exclude=G101,G304,G301,G306 ./...

govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	go test ./...

build:
	go generate ./...
	go build ./...

goreleaser:
	goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $(shell svu major)
tag-minor:
	git tag $(shell svu minor)
tag-patch:
	git tag $(shell svu patch)

release:
	git push origin --tags
	GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/discord-bot@$(shell svu current)

bump-glazed:
	@if [ -z "$(GO_GO_GOLEMS_MODULES)" ]; then \
		echo "No github.com/go-go-golems modules found in go.mod/go.sum"; \
		exit 1; \
	fi
	go get $(addsuffix @latest,$(GO_GO_GOLEMS_MODULES))
	go mod tidy

discord-bot_BINARY=$(shell which discord-bot)
install:
	go build -o ./dist/discord-bot ./cmd/discord-bot && \
		cp ./dist/discord-bot $(discord-bot_BINARY)
