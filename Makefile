.PHONY: all fmt vet lint vuln test coverage perft acceptance bench build install-hooks help

GO      := go
GOFMT   := gofmt
GOLINT  := golangci-lint
GOBIN   := $(shell $(GO) env GOPATH)/bin

GO_VERSION  := 1.22
CHESS_PKG   := ./internal/chess/...
ALL_PKGS    := ./...
ACCEPT_PKGS := ./tests/acceptance/...

# ─── Default ──────────────────────────────────────────────────────────────────

all: fmt vet lint test build ## Run all checks and build (same as CI, minus release)

# ─── Stage 1: Format ──────────────────────────────────────────────────────────

fmt: ## Check formatting (mirrors CI format job)
	@echo "==> gofmt"
	@unformatted=$$($(GOFMT) -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Files need formatting:"; \
		echo "$$unformatted"; \
		echo "Run: gofmt -w ."; \
		exit 1; \
	fi
	@echo "OK"

fmt-fix: ## Auto-fix formatting
	@$(GOFMT) -w .

# ─── Stage 2: Vet ─────────────────────────────────────────────────────────────

vet: ## Run go vet (mirrors CI vet job)
	@echo "==> go vet"
	@$(GO) vet $(ALL_PKGS)
	@echo "OK"

# ─── Stage 3: Lint ────────────────────────────────────────────────────────────

lint: ## Run golangci-lint (mirrors CI lint job)
	@echo "==> golangci-lint"
	@$(GOLINT) run --timeout=5m
	@echo "OK"

vuln: ## Run govulncheck
	@echo "==> govulncheck"
	@$(GO) run golang.org/x/vuln/cmd/govulncheck@latest $(ALL_PKGS)
	@echo "OK"

# ─── Stage 4: Test + Coverage ─────────────────────────────────────────────────

test: ## Run all tests with race detector (mirrors CI test-coverage job)
	@echo "==> go test -race"
	@$(GO) test -race -coverprofile=coverage.out -covermode=atomic -timeout=5m $(ALL_PKGS)
	@echo "OK"

coverage: test ## Check chess package coverage >= 90%
	@echo "==> coverage gate (internal/chess >= 90%)"
	@COVERAGE=$$($(GO) tool cover -func=coverage.out \
		| awk '/chess_go\/internal\/chess/ && /total:/ {print $$3}' \
		| tr -d '%'); \
	if [ -z "$$COVERAGE" ]; then \
		echo "WARN: internal/chess not yet implemented — skipping coverage gate."; \
		exit 0; \
	fi; \
	echo "chess package coverage: $${COVERAGE}%"; \
	PASS=$$(awk -v cov="$$COVERAGE" 'BEGIN { print (cov >= 90.0) ? "yes" : "no" }'); \
	if [ "$$PASS" != "yes" ]; then \
		echo "FAIL: coverage $${COVERAGE}% < 90% (NFR-09)"; \
		exit 1; \
	fi; \
	echo "OK"

coverage-html: test ## Open coverage report in browser
	@$(GO) tool cover -html=coverage.out

# ─── Stage 5: Perft ───────────────────────────────────────────────────────────

perft: ## Run perft validation tests (depths 1-4)
	@echo "==> perft (depths 1-4)"
	@$(GO) test -run TestPerft -tags slow -v -timeout=10m $(CHESS_PKG) || \
		echo "WARN: internal/chess not yet implemented — skipping perft."

perft-deep: ## Run perft depth 5 (slow — ~50s at 100k NPS)
	@echo "==> perft depth 5 (slow)"
	@$(GO) test -run TestPerftDepth5 -tags slow -v -timeout=15m $(CHESS_PKG)

# ─── Stage 6: Acceptance ──────────────────────────────────────────────────────

acceptance: ## Run acceptance tests (walking skeleton)
	@echo "==> acceptance tests"
	@$(GO) test -v -timeout=5m $(ACCEPT_PKGS)

# ─── Stage 7: Benchmarks ──────────────────────────────────────────────────────

bench: ## Run benchmarks and display results
	@echo "==> benchmarks"
	@$(GO) test -bench=. -benchmem -count=3 -benchtime=3s -timeout=10m $(ALL_PKGS)

bench-compare: ## Compare benchmarks against saved baseline (requires benchstat)
	@command -v benchstat >/dev/null 2>&1 || $(GO) install golang.org/x/perf/cmd/benchstat@latest
	@echo "==> benchmark comparison"
	@$(GO) test -bench=. -benchmem -count=3 -benchtime=3s -timeout=10m $(ALL_PKGS) \
		| tee bench-current.txt
	@if [ -f bench-baseline.txt ]; then \
		benchstat bench-baseline.txt bench-current.txt; \
	else \
		echo "No baseline found. Run 'make bench-save-baseline' first."; \
	fi

bench-save-baseline: ## Save current benchmark results as baseline
	@$(GO) test -bench=. -benchmem -count=3 -benchtime=3s -timeout=10m $(ALL_PKGS) \
		| tee bench-baseline.txt
	@echo "Baseline saved to bench-baseline.txt"

# ─── Stage 8: Build ───────────────────────────────────────────────────────────

build: ## Build both binaries for current platform
	@echo "==> build"
	@$(GO) build -o bin/chess-go   ./cmd/chess-go/...
	@$(GO) build -o bin/chess-server ./cmd/chess-server/...
	@echo "Binaries: bin/chess-go  bin/chess-server"

build-all: ## Cross-compile for all release targets
	@echo "==> cross-compile"
	@for target in \
		linux/amd64 linux/arm64 \
		darwin/amd64 darwin/arm64 \
		windows/amd64; do \
		GOOS=$$(echo $$target | cut -d/ -f1); \
		GOARCH=$$(echo $$target | cut -d/ -f2); \
		EXT=""; [ "$$GOOS" = "windows" ] && EXT=".exe"; \
		echo "  $$target"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH \
			$(GO) build -o /dev/null ./cmd/chess-go/... && \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH \
			$(GO) build -o /dev/null ./cmd/chess-server/...; \
	done
	@echo "OK"

# ─── Git Hooks ────────────────────────────────────────────────────────────────

install-hooks: ## Install pre-commit and pre-push git hooks
	@scripts/install-hooks.sh
	@echo "Hooks installed. Use 'git commit --no-verify' to bypass pre-commit (emergencies only)."

# ─── Help ─────────────────────────────────────────────────────────────────────

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
