.PHONY: build test clean run install-deps release release-test release-snapshot

# Build commands (no CGO required with modernc.org/sqlite)
build:
	mkdir -p build
	CGO_ENABLED=0 go build -o build/lazysmtp ./src

build-all:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/lazysmtp-linux-amd64 ./src
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/lazysmtp-darwin-amd64 ./src
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/lazysmtp-darwin-arm64 ./src
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/lazysmtp-windows-amd64.exe ./src

# Test commands
test:
	CGO_ENABLED=0 go test -v -cover ./...

test-race:
	CGO_ENABLED=0 go test -race -v ./...

bench:
	CGO_ENABLED=0 go test -bench=. -benchmem ./...

# Run commands
run: build
	./build/lazysmtp

# Clean commands
clean:
	rm -rf build/
	rm -f *.test
	rm -f coverage.out coverage.html
	rm -rf dist/

# Dependency commands
install-deps:
	go install github.com/goreleaser/goreleaser@latest

# Format and lint
fmt:
	go fmt ./...

vet:
	go vet ./...

# Release commands (CGO_ENABLED=0 for pure Go)
release:
	CGO_ENABLED=0 goreleaser release --clean

release-test:
	CGO_ENABLED=0 goreleaser release --skip=publish --clean

release-snapshot:
	CGO_ENABLED=0 goreleaser release --snapshot --clean

# Development commands
dev:
	go run ./src

# Install locally
install: build
	install -m 755 build/lazysmtp /usr/local/bin/lazysmtp

# Uninstall
uninstall:
	rm -f /usr/local/bin/lazysmtp

# Coverage
coverage:
	CGO_ENABLED=0 go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Show paths (for debugging)
show-paths:
	@echo "Config: $(shell go run -exec echo {{.ConfigPath}} 2>/dev/null || echo 'Run build first')"
	@echo "Data: $(shell go run -exec echo {{.DataPath}} 2>/dev/null || echo 'Run build first')"
