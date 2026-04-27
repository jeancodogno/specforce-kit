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

VERSION=$(shell node -p "require('./package.json').version")
PLATFORMS=linux-x64 linux-arm64 darwin-x64 darwin-arm64 win32-x64

npm-packages:
	@echo "Generating NPM native packages for version $(VERSION)..."
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d- -f1); \
		arch=$$(echo $$platform | cut -d- -f2); \
		goos=$$os; \
		goarch=$$arch; \
		[ "$$os" = "win32" ] && goos="windows"; \
		[ "$$arch" = "x64" ] && goarch="amd64"; \
		pkg_name="@jeancodogno/specforce-kit-$$os-$$arch"; \
		pkg_dir="npm/$$os-$$arch"; \
		mkdir -p $$pkg_dir/bin; \
		src_dir=$$(ls -d dist/specforce-kit_$${goos}_$${goarch}* | head -n 1); \
		binary_name="specforce"; \
		[ "$$os" = "win32" ] && binary_name="specforce.exe"; \
		cp $$src_dir/$$binary_name $$pkg_dir/bin/specforce$$( [ "$$os" = "win32" ] && echo ".exe" ); \
		cp LICENSE $$pkg_dir/LICENSE; \
		cp README.md $$pkg_dir/README.md; \
		printf '{\n  "name": "%s",\n  "version": "%s",\n  "description": "Native binary for %s %s",\n  "license": "MIT",\n  "author": "Jean Codogno <jeancarlo.eng.comp@gmail.com>",\n  "repository": { "type": "git", "url": "https://github.com/jeancodogno/specforce-kit.git" },\n  "bugs": { "url": "https://github.com/jeancodogno/specforce-kit/issues" },\n  "homepage": "https://github.com/jeancodogno/specforce-kit#readme",\n  "os": ["%s"],\n  "cpu": ["%s"],\n  "bin": { "specforce": "bin/specforce%s" }\n}\n' \
			"$$pkg_name" "$(VERSION)" "$$os" "$$arch" "$$os" "$$arch" "$$( [ "$$os" = "win32" ] && echo ".exe" )" > $$pkg_dir/package.json; \
	done

.PHONY: build test lint security quality check bootstrap clean-test install uninstall npm-packages
