.PHONY: all test build lint lintmax docker-lint gosec govulncheck goreleaser tag-major tag-minor tag-patch release bump-glazed install

all: test build

TAPES=$(shell ls doc/vhs/*tape)
gifs: $(TAPES)
	for i in $(TAPES); do vhs < $$i; done

docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run -v

lint:
	golangci-lint run -v # Basic lint run

lintmax: # New target for more comprehensive local linting
	golangci-lint run -v --max-same-issues=100

gosec: # New target for GoSec scan
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	# Adjust exclusions as needed, mirroring the workflow
	gosec -exclude=G101,G304,G301,G306,G204 -exclude-dir=.history ./...

govulncheck: # New target for govulncheck scan
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	go test ./...

build:
	go generate ./...
	go build ./...

goreleaser:
	goreleaser release --skip=sign --snapshot --clean

tag-major:
	git tag $(shell svu major)

tag-minor:
	git tag $(shell svu minor)

tag-patch:
	git tag $(shell svu patch)

release:
	git push --tags
	GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/go-emrichen@$(shell svu current)

bump-glazed:
	go get github.com/go-go-golems/glazed@latest
	go mod tidy

codeql-local:
	@echo "To run CodeQL locally, you need to install the CodeQL CLI from https://github.com/github/codeql-cli-binaries/releases"
	@echo "And clone the CodeQL queries from https://github.com/github/codeql-go.git"
	@if command -v codeql >/dev/null 2>&1; then \
		codeql database create --language=go --source-root=. ./codeql-db && \
		codeql database analyze ./codeql-db "$(HOME)/codeql-go/ql/src/go/Security" --format=sarif-latest --output=codeql-results.sarif && \
		echo "Results saved to codeql-results.sarif"; \
	else \
		echo "CodeQL CLI not found. Skipping analysis."; \
	fi

emrichen_BINARY=$(shell which emrichen)
install:
	go build -o ./dist/emrichen ./cmd/emrichen && \
		cp ./dist/emrichen $(emrichen_BINARY)
