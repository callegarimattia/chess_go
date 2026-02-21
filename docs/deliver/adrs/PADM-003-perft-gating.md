# PADM-003: Perft Test Gating Strategy

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Apex (platform-architect)
**Affected components**: `.github/workflows/ci.yml` (perft job), `internal/chess/perft_test.go`

---

## Context

Perft (performance test) is the canonical correctness oracle for chess move generators. It counts the exact number of leaf nodes at depth N from a given position. Any deviation from the known value indicates a bug in move generation (either illegal moves are generated, or legal moves are missing).

NFR-02 requires perft results to match known values at depth 5. The challenge is that perft depth 5 from the starting position = 4,865,609 nodes. At 100k NPS this takes ~49 seconds. Running depth 5 on every PR push would add nearly a minute to every CI run on a branch that has not yet implemented the engine.

---

## Decision

- **Standard CI** (every push, every PR): run perft at **depths 1-4** only. Depth 4 = 197,281 nodes, completing in ~2 seconds at 100k NPS.
- **Slow tests** (opt-in): depth 5 is gated behind the `-tags slow` build tag and is run only when explicitly invoked (e.g., `go test -run TestPerft -tags slow ./internal/chess/...`).
- **Release pipeline**: inherits the same depth 1-4 gate. Depth 5 validation is a developer responsibility before tagging.
- **Multiple positions**: Standard CI also runs Kiwipete at depths 1-3 to catch position-specific move generator bugs (castling, en passant, promotion).

---

## Alternatives Considered

### Alternative A: Run depth 5 on every push

**Evaluation**: Maximum correctness confidence. Adds ~50 seconds to every CI run.

**Rejection**: At 100k NPS, depth 5 is 49 seconds. This adds nearly a minute to every branch push and PR update. The incremental correctness gain over depth 4 (197k vs 4.8M nodes) does not justify the latency cost for the development feedback loop. Depth 4 catches the vast majority of move generator bugs; depth 5 catches only exotic combinations that require deep position repetition.

### Alternative B: Skip perft in CI; run only locally

**Evaluation**: Zero CI time cost. Maximum developer flexibility.

**Rejection**: Correctness is a Critical quality attribute (architecture-design.md, Section 6). Removing the correctness gate from CI means a broken move generator can merge to main. The perft gate must be automated.

### Alternative C: Run perft only on `main` pushes, not PRs

**Evaluation**: Reduces PR cycle time; still validates before main merge.

**Rejection**: Correctness gates must block merges, not just validate after. A move generator bug should be caught before the PR is merged, not after. The perft gate must run on PRs.

---

## Consequences

**Positive**:
- Depth 1-4 completes in ~2-4 seconds; negligible CI overhead
- Kiwipete coverage catches castling, en passant, and promotion edge cases in standard CI
- `-tags slow` allows depth 5 validation without mandating it in every CI run
- Build tag approach is idiomatic Go; no external test filtering configuration required

**Negative**:
- Depth 5 is not guaranteed to have passed before a release tag is created; this is a documented developer responsibility
- A move generator bug that manifests only at depth 5 could reach a release tag; mitigated by the developer running `go test -tags slow` locally before tagging
- Position coverage in standard CI is limited to 2 positions (starting + Kiwipete); additional positions (Position 3, 4, 5 from CPW) are in the slow suite only
