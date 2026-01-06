# lazySMTP Commands

This document lists all available commands for lazySMTP.

## Make Commands

### Build Commands

```bash
make build
```
Build binary for current platform (CGO disabled).
Output: `build/lazysmtp`

```bash
make build-all
```
Build binaries for all supported platforms to `build/` directory:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Output files:
- `build/lazysmtp` (current platform)
- `build/lazysmtp-linux-amd64`
- `build/lazysmtp-darwin-amd64`
- `build/lazysmtp-darwin-arm64`
- `build/lazysmtp-windows-amd64.exe`

### Test Commands

```bash
make test
```
Run all tests with coverage reporting.

```bash
make test-race
```
Run tests with race detector enabled.

```bash
make bench
```
Run benchmarks and show memory statistics.

### Run Commands

```bash
make run
```
Build and run lazysmtp in development mode.

```bash
make dev
```
Run lazysmtp directly using `go run` (for rapid iteration).

### Clean Commands

```bash
make clean
```
Remove all generated files:
- Build directory (`build/`)
- Test files (`*.test`)
- Coverage files (`coverage.out`, `coverage.html`)
- Distribution directory (`dist/`)

### Dependency Commands

```bash
make install-deps
```
Install goreleaser for release management:
```bash
go install github.com/goreleaser/goreleaser@latest
```

### Code Quality

```bash
make fmt
```
Format all Go code using `go fmt`.

```bash
make vet
```
Run `go vet` to check for common mistakes.

### Release Commands

```bash
make release
```
Create a full release:
- Build binaries for all platforms
- Generate checksums
- Create GitHub release
- Publish to AUR
- Generate packages (deb, rpm, apk)

**Environment variables required:**
- `GITHUB_TOKEN`: GitHub personal access token
- `AUR_PRIVATE_KEY`: Path to AUR SSH private key

```bash
make release-test
```
Test release workflow without publishing:
- Builds everything
- Skips GitHub publish
- Skips AUR upload

```bash
make release-snapshot
```
Create a snapshot release for testing:
- Doesn't tag version
- Generates artifacts in `dist/`

### Installation Commands

```bash
make install
```
Install lazysmtp from `build/lazysmtp` to `/usr/local/bin/lazysmtp`

```bash
make uninstall
```
Remove lazysmtp from `/usr/local/bin/lazysmtp`

### Coverage Commands

```bash
make coverage
```
Generate test coverage report:
- Creates `coverage.out` file
- Generates `coverage.html` for viewing in browser

## Binary Commands

```bash
lazysmtp
```
Start lazySMTP with default settings:
- Port: 2525
- Database: XDG data directory

```bash
lazysmtp -port 1025
```
Start on custom port (useful when running as non-root).

```bash
lazysmtp -db /path/to/custom.db
```
Use a custom database file location.

```bash
lazysmtp -h
lazysmtp --help
```
Show help message with all options.

## TUI Commands (Inside Application)

### Navigation

- `j` - Move down in email list
- `k` - Move up in email list
- `d` - Delete selected email
- `SPACE` - Toggle SMTP server on/off
- `q` - Quit application

## Go Commands

```bash
go run ./src
```
Run lazysmtp without building (for development).

```bash
go build -o build/lazysmtp ./src
```
Manual build to build directory (same as `make build`).

```bash
go test ./src -v
```
Run tests with verbose output.

```bash
go test ./src -cover
```
Run tests with coverage.

```bash
go test ./src -bench=. -benchmem
```
Run benchmarks.

```bash
go fmt ./src
```
Format Go code.

```bash
go vet ./src
```
Check for code issues.

```bash
go mod tidy
```
Clean up go.mod dependencies.

```bash
go mod download
```
Download all dependencies.

## Goreleaser Commands

```bash
goreleaser check
```
Validate `.goreleaser.yml` configuration.

```bash
goreleaser release --clean
```
Full release (same as `make release`).

```bash
goreleaser release --skip=publish --clean
```
Test release (same as `make release-test`).

```bash
goreleaser release --snapshot --clean
```
Snapshot release (same as `make release-snapshot`).

## Development Workflow

### Typical Development Cycle

```bash
# 1. Make changes
vim src/main.go

# 2. Run tests
make test

# 3. Format code
make fmt

# 4. Check for issues
make vet

# 5. Build and test
make dev

# 6. Clean old artifacts
make clean
```

### Release Workflow

```bash
# 1. Update version in src/main.go
vim src/main.go

# 2. Commit and push changes
git add .
git commit -m "Bump version to v1.0.0"
git push

# 3. Tag release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. Test release locally
make release-snapshot

# 5. Dry run
make release-test

# 6. Actual release
make release
```

### AUR Release Workflow

```bash
# 1. Set up environment
export AUR_PRIVATE_KEY=~/.ssh/aur
export GITHUB_TOKEN=your_token

# 2. Ensure AUR package exists
git clone ssh://aur@aur.archlinux.org/lazysmtp-git.git

# 3. Release (goreleaser will update AUR)
make release
```

## Troubleshooting

### Build Fails with CGO Errors

```bash
# Ensure CGO is disabled
export CGO_ENABLED=0
make build
```

### Tests Fail with Database Errors

```bash
# Clean and rebuild
make clean
make build
make test
```

### Goreleaser Fails

```bash
# Check configuration
goreleaser check

# Try snapshot first
make release-snapshot
```

### AUR Upload Fails

```bash
# Verify SSH key works
ssh -T aur@aur.archlinux.org

# Check AUR package manually
cd lazysmtp-git
makepkg -si
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `make build` | Build binary |
| `make test` | Run tests |
| `make run` | Run app |
| `make clean` | Clean artifacts |
| `make release` | Create release |
| `lazysmtp` | Start app |
| `lazysmtp -port X` | Custom port |
| `lazysmtp -db X` | Custom database |
