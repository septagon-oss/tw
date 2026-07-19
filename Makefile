SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

GO ?= go
GOTMPDIR ?= $(CURDIR)/.tmp-go-tmp
GO_ENV := GOTMPDIR="$(GOTMPDIR)" TMPDIR="$(GOTMPDIR)"
GO_FILES := $(shell find . -type f -name '*.go' -print | sort)

.PHONY: fmt-check test vet verify

$(GOTMPDIR):
	mkdir -p "$@"

fmt-check:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [[ -n "$$unformatted" ]]; then \
		echo "Go files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

test: | $(GOTMPDIR)
	$(GO_ENV) $(GO) test ./...

vet: | $(GOTMPDIR)
	$(GO_ENV) $(GO) vet ./...

verify: fmt-check vet test
