# PADM-001: CI/CD Platform Selection

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Apex (platform-architect)
**Affected components**: `.github/workflows/ci.yml`, `.github/workflows/release.yml`

---

## Context

The project needs a CI/CD platform to enforce quality gates on every push and to automate binary release packaging. The platform must be free for a public Go repository, require minimal operational overhead, and integrate natively with GitHub (where the repository is hosted).

The deployment target is local developer machines. There is no cloud infrastructure to provision, no containers to build, and no servers to deploy to. The CI platform's responsibilities are limited to:
1. Running quality gates (format, vet, lint, test, perft, benchmark, build)
2. Building and publishing release binaries

---

## Decision

Use **GitHub Actions** with hosted runners (`ubuntu-latest`).

---

## Alternatives Considered

### Alternative A: CircleCI

**Evaluation**: Free tier limited to 6,000 build minutes/month; requires separate account and configuration. Native GitHub integration is weaker than GitHub Actions. No advantage for a public Go repo with simple pipeline requirements.

**Rejection**: GitHub Actions is specified in the deployment configuration. CircleCI adds account management overhead with no capability advantage for this pipeline.

### Alternative B: Local Makefile only (no CI)

**Evaluation**: A `Makefile` can run all quality checks locally. Zero infrastructure cost. No pipeline YAML to maintain.

**Rejection**: A Makefile cannot enforce gates on PRs or certify that the tagged commit passed all gates before release. "main is always releasable" requires automated enforcement on every push. The Makefile is a useful local developer tool but is insufficient as a sole quality gate.

### Alternative C: Self-hosted GitHub Actions runner

**Evaluation**: Faster runners, no per-minute billing, can use local machine resources.

**Rejection**: Self-hosted runners require infrastructure management and are unavailable for public repos receiving PRs from external contributors. GitHub-hosted runners are sufficient for the pipeline (estimated ~12 minutes per run). Self-hosted is a v2 optimization if runner cost becomes significant.

---

## Consequences

**Positive**:
- Zero infrastructure to operate; GitHub manages runner lifecycle
- Native PR status checks integration; branch protection rules consume job status directly
- Free unlimited minutes for public repositories
- `ubuntu-latest` provides Go toolchain, standard build tools, and GitHub CLI pre-installed
- Artifact storage and cache included

**Negative**:
- `ubuntu-latest` may change the underlying OS version; pin to `ubuntu-24.04` if stability is required
- GitHub Actions rate limits apply on private repos (not a concern for this public repo)
- Benchmarks on shared runners have variable latency; benchstat's statistical testing mitigates false regressions
