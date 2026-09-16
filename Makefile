# Use this module's pinned APIs for hooks and CI, even inside a multi-repo workspace.
# Explicit development overrides remain possible with `make GOWORK=...`.
export GOWORK := off
# Local validation needs no VCS stamp (also works in linked/restricted worktrees).
export GOFLAGS ?= -buildvcs=false

.PHONY: all test build lint lintmax golangci-lint-install gosec govulncheck \
        glazed-lint-build glazed-lint goreleaser tag-major tag-minor tag-patch release bump-glazed install

all: test build

VERSION=v0.0.1
GLAZED_LINT_BIN ?= $(CURDIR)/.bin/glazed-lint
GLAZED_LINT_PKG ?= github.com/go-go-golems/glazed/cmd/tools/glazed-lint
GLAZED_VERSION ?= $(shell GOWORK=off go list -m -f '{{.Version}}' github.com/go-go-golems/glazed 2>/dev/null)
GLAZED_LINT_FLAGS ?=
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
	@test -n "$(GLAZED_VERSION)" || (echo "Cannot resolve pinned Glazed module version"; exit 1)
	mkdir -p $(dir $(GLAZED_LINT_BIN))
	GOBIN=$(dir $(GLAZED_LINT_BIN)) go install $(GLAZED_LINT_PKG)@$(GLAZED_VERSION)

glazed-lint: glazed-lint-build
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)


lint: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) config verify
	$(GOLANGCI_LINT_BIN) run -v $(GOLANGCI_LINT_ARGS)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

lintmax: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) run -v --max-same-issues=100 $(GOLANGCI_LINT_ARGS)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

GOSEC_VERSION ?= v2.29.0
GOVULNCHECK_VERSION ?= v1.8.0
GOSEC_BIN ?= $(CURDIR)/.bin/gosec
GOVULNCHECK_BIN ?= $(CURDIR)/.bin/govulncheck

gosec:
	mkdir -p $(dir $(GOSEC_BIN))
	GOBIN=$(dir $(GOSEC_BIN)) go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	$(GOSEC_BIN) -exclude-generated -exclude=G101,G304,G301,G306,G204 -exclude-dir=.history -exclude-dir=ttmp ./...

govulncheck:
	mkdir -p $(dir $(GOVULNCHECK_BIN))
	GOBIN=$(dir $(GOVULNCHECK_BIN)) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	$(GOVULNCHECK_BIN) ./...

.PHONY: check hooks-install vet
check: lint test build vet logcopter-check gosec govulncheck

hooks-install:
	lefthook install

vet:
	go vet ./...

test:
	go test ./...

build:
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

.PHONY: bump-go-go-golems
bump-go-go-golems:
	@deps="$$(awk '/^require[[:space:]]+github\.com\/go-go-golems\// { print $$2 } /^[[:space:]]*github\.com\/go-go-golems\// { print $$1 }' go.mod | sort -u)"; \
	if [ -z "$$deps" ]; then \
		echo "No github.com/go-go-golems dependencies in go.mod"; \
	else \
		echo "Bumping go-go-golems dependencies:"; \
		echo "$$deps"; \
		for dep in $$deps; do GOWORK=off go get "$${dep}@latest"; done; \
	fi
	GOWORK=off go mod tidy

.PHONY: logcopter-generate
logcopter-generate:
	GOWORK=off go generate ./...

.PHONY: logcopter-check
logcopter-check:
	GOWORK=off go tool logcopter-gen -area-prefix go-go-golems.discord-bot -strip-prefix github.com/go-go-golems/discord-bot -check ./internal/... ./pkg/... ./cmd/...
