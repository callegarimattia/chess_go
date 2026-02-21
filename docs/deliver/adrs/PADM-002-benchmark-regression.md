# PADM-002: Benchmark Regression Detection Strategy

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Apex (platform-architect)
**Affected components**: `.github/workflows/ci.yml` (benchmarks job)

---

## Context

NFR-01 requires the engine to achieve > 100,000 nodes per second. The pipeline must detect regressions in search performance on every PR. Three challenges:

1. **CI runner variability**: GitHub-hosted runners share hardware; a single benchmark run can vary 15–30% due to scheduling noise.
2. **Baseline persistence**: The comparison target must survive across runs without a database.
3. **Threshold selection**: A fixed percentage threshold (e.g., 20%) must avoid false positives from noise while catching real regressions.

---

## Decision

Use **`benchstat` from `golang.org/x/perf`** to compare benchmark results between the current PR and a cached baseline from the `main` branch.

- Run benchmarks with `-count=3 -benchtime=3s` to produce multiple samples for statistical significance testing.
- Cache `bench-baseline.txt` in GitHub Actions cache, keyed by Go version and OS.
- Update the baseline only on successful pushes to `main`.
- Fail the pipeline if benchstat reports a regression > 20% with p < 0.05.
- On PRs, compare against the main branch baseline; on main pushes, the current run becomes the new baseline.

---

## Alternatives Considered

### Alternative A: Single-run percentage comparison

**Description**: Run benchmark once; compare the numeric value against a hardcoded constant (e.g., 100,000 NPS).

**Evaluation**: Simple but brittle. A single run on a noisy CI runner can vary 20–30%, producing false failures. A hardcoded constant is not sensitive to gradual regressions across multiple merges.

**Rejection**: Statistical comparison against a recent baseline is more accurate and more sensitive to real regressions. `benchstat` is purpose-built for this problem.

### Alternative B: External benchmark tracking service (Benchmarks.io, custom InfluxDB)

**Description**: Push benchmark results to an external service; visualise trends; alert on anomalies.

**Evaluation**: Provides richer trend analysis and per-commit history. Requires external infrastructure and integration.

**Rejection**: Infrastructure overhead is disproportionate for a solo/small team project. GitHub Actions cache provides sufficient baseline storage for a single-branch comparison. An external service becomes appropriate if the engine grows a community of contributors with competing performance claims.

### Alternative C: No automated regression detection; manual review

**Description**: Developer reviews benchmark output in CI logs manually; no automated failure.

**Evaluation**: Zero tooling cost. Relies on discipline to catch regressions.

**Rejection**: Manual review is not a quality gate. Silent regressions accumulate without automated detection. The 20% NPS degradation required before a user notices is achievable across a series of small undetected regressions.

---

## Consequences

**Positive**:
- `benchstat` applies Welch's t-test; filters out noise at the statistical level
- GitHub Actions cache is free and sufficient for storing a single baseline file
- `-count=3` with 3-second benchtime provides enough samples for significance testing
- The 20% threshold is generous enough to tolerate runner noise while catching real regressions

**Negative**:
- First run has no baseline; regression check is skipped (acceptable; baseline established on first main push)
- Benchmark results on shared runners still vary; a legitimately fast PR may show no statistically significant improvement even if the local machine shows one
- Cache eviction (after 7 days of no access) loses the baseline; next run skips comparison and establishes a new one
