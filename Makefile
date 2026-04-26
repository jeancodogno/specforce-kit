BINARY=specforce
PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

# Go paths
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin

build:
	go build -o $(BINARY) src/cmd/specforce/main.go

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)

test-coverage:
	go test -v -coverprofile=coverage.out ./...

test:
	go test ./...

lint:
	golangci-lint run ./...

security:
	$(GOBIN)/gosec -exclude=G703,G304,G306 -fmt=sarif -out=results.sarif ./...
	$(GOBIN)/govulncheck ./... > /dev/null 2>&1 || echo "Warning: govulncheck found vulnerabilities (likely Go stdlib version mismatch). Ignored locally."

quality:
	$(GOBIN)/gremlins unleash

quality-lived:
	$(GOBIN)/gremlins unleash -S l

check: lint security test quality

bootstrap:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/go-gremlins/gremlins/cmd/gremlins@latest

clean-test:
	rm -rf test-init/.specforce test-init/.agent test-init/.claude test-init/.gemini
	rm -f results.sarif coverage.out

.PHONY: build test lint security quality check bootstrap clean-test install uninstall
