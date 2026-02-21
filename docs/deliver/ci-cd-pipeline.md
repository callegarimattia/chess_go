# CI/CD Pipeline Design — chess-engine
**Feature**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved
**Author**: Apex (nw-platform-architect)

---

## 1. Pipeline Philosophy

The pipeline enforces a single invariant: **main is always releasable**. Every quality gate is a binary pass/fail decision. There are no warnings that are allowed to accumulate.

The pipeline is split into two workflows:
- `ci.yml` — runs on every push to `main` and every PR targeting `main`
- `release.yml` — runs only on `git tag v*`; extends CI with release binary packaging

---

## 2. CI Pipeline — Stage Design

### Stage Execution Model

Stages 1–3 (format, vet, lint) run in parallel as fast feedback. They are independent of one another and share no artifacts.

Stages 4–8 run sequentially after stages 1–3 pass. Stage 4 (test+coverage) produces `coverage.out` consumed by the coverage threshold check. Stage 7 (benchmarks) consumes a cached baseline from a prior run.

```
┌──────────┐  ┌───────┐  ┌────────┐
│  format  │  │  vet  │  │  lint  │  ← parallel, fast (~60s total)
└────┬─────┘  └───┬───┘  └───┬────┘
     └────────────┴──────────┘
                  │ all pass
                  ▼
         ┌────────────────┐
         │  test+coverage │  ← go test -race, coverage threshold (~3-5min)
         └───────┬────────┘
                 │
                 ▼
         ┌───────────────┐
         │     perft     │  ← depth 1-4 validation (~2-3min)
         └───────┬───────┘
                 │
                 ▼
         ┌───────────────┐
         │  acceptance   │  ← walking skeleton test (~30s)
         └───────┬───────┘
                 │
                 ▼
         ┌───────────────┐
         │  benchmarks   │  ← NPS regression check (~2min)
         └───────┬───────┘
                 │
                 ▼
         ┌───────────────┐
         │  build (3x)   │  ← linux/darwin/windows cross-compile (~1min)
         └───────────────┘
```

**Total estimated CI time: ~10-12 minutes** (dominated by perft at depth 4 and race-enabled tests).

---

## 3. Stage Specifications

### Stage 1: Format Check
**Command**: `gofmt -l .`
**Pass criteria**: zero output (no files need formatting)
**Fail action**: print list of unformatted files, exit non-zero
**Rationale**: Enforces consistent formatting without style debates. `gofmt` is canonical Go style.

### Stage 2: Vet
**Command**: `go vet ./...`
**Pass criteria**: exit code 0
**Fail action**: print vet findings, exit non-zero
**Rationale**: Catches common bugs (unreachable code, misused format verbs, suspicious constructs) before tests run.

### Stage 3: Lint
**Command**: `golangci-lint run --timeout=5m`
**Linters enabled**: `staticcheck`, `errcheck`, `unused`, `govet` (built-in), `ineffassign`, `gosimple`
**Config file**: `.golangci.yml` in project root (managed separately)
**Pass criteria**: zero issues
**Rationale**: staticcheck catches semantic bugs; errcheck prevents ignored errors; unused catches dead code. All three are specified in technology-stack.md.

### Stage 4: Unit Tests + Coverage
**Command**:
```bash
go test -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | grep "chess_go/internal/chess" | tail -1
```
**Pass criteria**:
- All tests pass
- `internal/chess` package coverage >= 90% (NFR-09)
- No race conditions detected
**Coverage extraction**: Parse `go tool cover -func` output for `internal/chess` total line. Fail if < 90.0%.
**Artifact**: `coverage.out` uploaded as workflow artifact for inspection.
**Rationale**: Race detector is mandatory (ADR-004); chess package coverage is a hard NFR.

### Stage 5: Perft Validation
**Command**: `go test -run TestPerft -tags slow -v -timeout 10m ./internal/chess/...`
**Pass criteria**: All perft assertions pass at depths 1–4 against known values (depth 5 is slow-tagged and skipped in standard CI; run only with explicit `-tags slow`):

| Depth | Node Count (starting position) |
|---|---|
| 1 | 20 |
| 2 | 400 |
| 3 | 8,902 |
| 4 | 197,281 |

Additional positions (Kiwipete, Position 3, 4, 5) tested at depth 1-3.
**Rationale**: Perft is the canonical correctness oracle for move generators. A wrong perft count means illegal moves are generated or legal moves are missing. This gate cannot be skipped.

**Performance note**: Perft depth 4 = 197,281 nodes. At 100k NPS this completes in ~2 seconds. The 10-minute timeout is a safety margin for slow CI runners.

### Stage 6: Acceptance Tests
**Command**: `go test -v -timeout 5m ./tests/acceptance/...`
**Pass criteria**: All non-`@skip` tests pass; walking skeleton scenarios pass.
**Rationale**: Validates that all architectural layers are wired together (chess + engine + tui). Walking skeleton is the integration proof.

### Stage 7: Benchmark Baseline
**Command**:
```bash
go test -bench=. -benchmem -count=3 -benchtime=3s ./... | tee bench-current.txt
```
**Regression check**: Compare against `bench-baseline.txt` cached in GitHub Actions cache. If baseline exists, run `benchstat bench-baseline.txt bench-current.txt` and fail if any benchmark shows a statistically significant regression >20%.

**Baseline management**:
- On `main` branch pushes that pass all prior stages, current results become the new baseline (cached with key `bench-baseline-{go-version}`)
- On PR branches, comparison is against the main branch baseline; regression fails the PR
- NPS target (NFR-01): `BenchmarkSearch` in `internal/engine/bench_test.go` must report >= 100,000 nodes/second

**Critical benchmarks**:

| Benchmark | Location | Target |
|---|---|---|
| `BenchmarkSearchDepth4` | `internal/engine/bench_test.go` | NPS >= 100k |
| `BenchmarkSearchDepth5` | `internal/engine/bench_test.go` | NPS >= 100k |
| `BenchmarkSearchDepth6` | `internal/engine/bench_test.go` | NPS >= 100k |
| `BenchmarkLegalMoves` | `internal/chess/...` | < 1µs/op |

### Stage 8: Cross-Platform Build
**Commands** (run in parallel as separate steps):
```bash
GOOS=linux   GOARCH=amd64 go build -o /dev/null ./cmd/...
GOOS=linux   GOARCH=arm64 go build -o /dev/null ./cmd/...
GOOS=darwin  GOARCH=amd64 go build -o /dev/null ./cmd/...
GOOS=darwin  GOARCH=arm64 go build -o /dev/null ./cmd/...
GOOS=windows GOARCH=amd64 go build -o /dev/null ./cmd/...
```
**Pass criteria**: All 5 × 2 builds (10 total: chess-go + chess-server per target) exit code 0.
**Rationale**: NFR — works on Linux, macOS, Windows. Cross-compilation is free in Go; this gate catches platform-specific stdlib usage.

---

## 4. Release Pipeline

### Trigger
`on: push: tags: ['v*']`

### Steps
1. **Checkout** with full tag history
2. **Run full CI suite** (identical to ci.yml stages 1-8)
3. **Build release binaries** for all 5 targets:
   - `chess-go-linux-amd64`
   - `chess-go-linux-arm64`
   - `chess-go-darwin-amd64`
   - `chess-go-darwin-arm64`
   - `chess-go-windows-amd64.exe`
   - `chess-server-linux-amd64`
   - `chess-server-linux-arm64`
   - `chess-server-darwin-amd64`
   - `chess-server-darwin-arm64`
   - `chess-server-windows-amd64.exe`
4. **Create GitHub release** with:
   - All 10 binary assets attached
   - Release notes auto-generated from Conventional Commits since the prior tag
   - Pre-release flag set for `v0.x.0` tags

### Build flags
```bash
go build -ldflags="-s -w -X main.version=${TAG_NAME}" -trimpath
```
`-s -w` strips debug info (~30% binary size reduction). `-trimpath` removes local filesystem paths from binary. `-X main.version` injects the tag name for future `--version` flag support.

---

## 5. Pipeline Security

| Concern | Implementation |
|---|---|
| Secrets | Only `GITHUB_TOKEN` (auto-provided); no additional secrets |
| Permissions | `contents: write` on release job only; CI job runs with minimal `read` |
| Dependency integrity | `go.sum` committed; Go module proxy verification |
| SAST | `go vet` + `staticcheck` (via golangci-lint) |
| Vulnerability scan | `govulncheck ./...` in lint stage |
| Runner | `ubuntu-latest` (GitHub-hosted); no self-hosted runners needed |

---

## 6. Failure Handling and Notifications

| Failure | Action |
|---|---|
| Format/vet/lint failure | PR blocked; developer fixes and re-pushes |
| Test failure | PR blocked; investigate test output in Actions log |
| Perft mismatch | PR blocked; CRITICAL — move generator has a correctness bug |
| Coverage below 90% | PR blocked; add tests before merging |
| Benchmark regression >20% | PR blocked; investigate and optimise or update baseline |
| Build failure on any target | PR blocked; fix platform-specific code |
| Release workflow failure | Tag is orphaned; delete tag, fix issue, re-tag |

GitHub branch protection rules enforce that CI must pass before merge. This is configured in the repository settings (documented in `branching-strategy.md`).

---

## 7. Caching Strategy

| Cache Key | Contents | Invalidation |
|---|---|---|
| `go-mod-{hash(go.sum)}` | Go module download cache | On go.sum change |
| `golangci-lint-{version}` | golangci-lint binary | On version bump in workflow |
| `bench-baseline-{go-version}` | `bench-baseline.txt` | Replaced on each main branch push |

Go module cache hit reduces stage 4-8 time by ~60% after the first run.
