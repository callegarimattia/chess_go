# Platform Architecture — chess-engine
**Feature**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved
**Author**: Apex (nw-platform-architect)

---

## 1. Platform Summary

This document describes the delivery infrastructure for the chess-engine feature. The platform is intentionally minimal: the product is a pair of static Go binaries distributed directly to developer machines. There is no cloud, no containers, no runtime infrastructure. The platform's job is exclusively CI quality enforcement and binary release packaging.

### Guiding Constraints

| Constraint | Source | Implication |
|---|---|---|
| Deployment target is local developer machine | Configuration | No cloud provider, no container runtime, no load balancer |
| No container orchestration | Configuration | No Docker, no Kubernetes, no registry |
| CI platform: GitHub Actions | Configuration | YAML workflows in `.github/workflows/` |
| Greenfield project | Configuration | Design everything from scratch; no existing CI to reuse |
| Observability: lightweight | Configuration | Structured stderr logs in binaries + Go benchmark suite |
| Deployment strategy: recreate | Configuration | Users download a new binary; no rolling or canary required |
| Branching: Trunk-Based Development | Configuration | Single main branch, short-lived feature branches <1 day |

---

## 2. Platform Components

### 2.1 Rejected Simpler Alternatives (documented per Principle 4)

Before finalising this design, two simpler alternatives were evaluated:

**Alternative 1: Makefile only, no CI**
Rejected. Makefile can run local checks but cannot enforce quality gates on every PR or produce verified release binaries. A solo developer skipping a local check before pushing breaks the "main is always releasable" guarantee. GitHub Actions is required.

**Alternative 2: Single monolithic workflow file**
Rejected. A single workflow cannot separate concerns between CI (runs on every push/PR) and release (runs on tag). Merging them creates a workflow that is harder to reason about and wastes release-build time on non-release events. Two files (ci.yml, release.yml) is the simplest correct split.

### 2.2 Chosen Infrastructure Components

| Component | Technology | Justification |
|---|---|---|
| CI pipeline | GitHub Actions | Specified; free for public repos; native GitHub integration |
| Lint | golangci-lint | Aggregates staticcheck, errcheck; specified in technology-stack |
| Coverage enforcement | `go test -coverprofile` + shell threshold check | stdlib; no external coverage service needed |
| Perft validation | `go test -run TestPerft` with `-tags slow` | stdlib test runner; known correctness gate |
| Benchmark baseline | `go test -bench=. -benchmem` + benchstat | benchstat is the standard Go tool for regression detection |
| Cross-compilation | `GOOS/GOARCH go build` | stdlib; Go's built-in cross-compilation |
| Release packaging | `gh release create` via `softprops/action-gh-release` | Standard GitHub Actions release action |
| Changelog | Auto-generated from Conventional Commits via GitHub's `--generate-notes` | No additional tooling required |

### 2.3 Components Not Used (and Why)

| Technology | Reason Not Used |
|---|---|
| Docker | No containers; binaries are deployed directly |
| Kubernetes | No container orchestration; overkill for local tool |
| Terraform / Pulumi | No cloud infrastructure to manage |
| Datadog / Prometheus | No long-running server in production context; local stderr logs sufficient |
| Dependabot / Renovate | Only one optional external dep (`golang.org/x/term`); low maintenance burden |
| Code coverage services (Codecov, Coveralls) | Coverage threshold enforced directly in CI; external service adds no value for greenfield solo project |

---

## 3. Delivery Pipeline Overview

```
Every push to main / every PR
────────────────────────────────────────────────────────────────
Stage 1: Format          gofmt -l .
Stage 2: Vet             go vet ./...
Stage 3: Lint            golangci-lint run
Stage 4: Test+Coverage   go test -race -coverprofile=coverage.out ./...
                         → fail if chess package coverage < 90%
Stage 5: Perft           go test -run TestPerft -tags slow ./internal/chess/...
                         → fail if depth 1-4 node counts mismatch known values
Stage 6: Acceptance      go test ./tests/acceptance/...
                         → walking skeleton must pass
Stage 7: Benchmarks      go test -bench=. -benchmem -count=3 ./...
                         → store results; compare vs baseline; fail if NPS regression > 20%
Stage 8: Build           GOOS=linux   GOARCH=amd64 go build ./cmd/...
                         GOOS=darwin  GOARCH=arm64 go build ./cmd/...
                         GOOS=windows GOARCH=amd64 go build ./cmd/...

On git tag v*
────────────────────────────────────────────────────────────────
Step 1: Run full CI suite (stages 1-8 above)
Step 2: Build release binaries for 5 targets
Step 3: Create GitHub release with binary assets + auto-changelog
```

---

## 4. Deployment Architecture (Local Binary Distribution)

Since deployment is local binary replacement, there is no deployment infrastructure to operate:

```
Release workflow produces:
  chess-go-linux-amd64
  chess-go-linux-arm64
  chess-go-darwin-amd64
  chess-go-darwin-arm64
  chess-go-windows-amd64.exe
  chess-server-linux-amd64
  chess-server-linux-arm64
  chess-server-darwin-amd64
  chess-server-darwin-arm64
  chess-server-windows-amd64.exe

User deployment procedure:
  1. Download binary for their OS/arch from GitHub Releases
  2. Replace existing binary on PATH (or first install)
  3. Verify: chess-go --version (future) or run and observe
```

**Rollback procedure** (required per Principle 7 — Rollback-First):
Since users hold binaries locally, rollback is:
1. User downloads the previous release tag binary from GitHub Releases
2. Replaces current binary
3. Verification: run binary, confirm prior behaviour

GitHub Releases retains all prior tags permanently. Every release is a rollback point.

---

## 5. DORA Metrics Targets

| Metric | Current (greenfield) | Target |
|---|---|---|
| Deployment Frequency | N/A | Each merged PR to main is a potential release; tag as needed |
| Lead Time for Changes | N/A | < 4 hours (commit to tagged release) |
| Change Failure Rate | N/A | < 10% (CI gates prevent most failures reaching a tag) |
| Time to Restore | N/A | < 1 hour (download previous release binary) |

These targets place the project in the **Medium performer** band (Accelerate). The CI pipeline enforces quality gates that directly reduce change failure rate.

---

## 6. Security Posture

For a local CLI tool with no network surface and no user data, the security model is:

| Concern | Mitigation |
|---|---|
| Supply chain | Go modules with checksums (`go.sum`); minimal external deps |
| Secret leakage | No secrets in source; GitHub Actions secrets for GITHUB_TOKEN only |
| SAST | `go vet` + `staticcheck` (via golangci-lint) on every PR |
| Race conditions | `-race` flag on all test runs |
| Dependency vulnerabilities | `govulncheck ./...` in CI (lightweight SCA) |

No DAST or SBOM required for a local binary tool with no network service in CI.

---

## 7. Platform ADRs

See `docs/deliver/adrs/` for platform-level architecture decisions:
- PADM-001: GitHub Actions over alternatives
- PADM-002: Benchmark regression via benchstat
- PADM-003: Perft test gating strategy
