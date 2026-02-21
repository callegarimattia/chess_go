# ADR-002: Search Algorithm

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Morgan (solution architect)
**Affected components**: `internal/engine` (search.go, eval.go)

---

## Context

The engine package requires a search algorithm to select the best move from a given position within a time budget. The algorithm must:
- Return a legal move within `movetime + 50ms` (NFR-03 / AC-13)
- Search at least depth 3 in 100ms (AC-12-02)
- Find forced mate in 1 at depth >= 1 (AC-12-03)
- Emit UCI info lines with depth, score, nodes, NPS, and PV (FR-02-06)
- Achieve > 100,000 NPS on modern hardware (NFR-01)

The algorithm design is specified in requirements as "alpha-beta with iterative deepening" (FR-02-01), with quiescence search (FR-02-07) and move ordering (FR-02-08) as Should-priority enhancements.

The architecture decision here concerns the specific implementation choices within the alpha-beta family: iterative deepening vs fixed-depth, fail-soft vs fail-hard, and the concurrency model for time management.

---

## Decision

Implement **iterative deepening alpha-beta (IDAS)** with:
- Fail-soft alpha-beta (returns exact score at root, enables better move ordering)
- Negamax frame (single recursive function handles both colors via score negation)
- Iterative deepening: search depth 1, then 2, then 3... until time expires; return best move from last completed depth
- Move ordering: MVV-LVA (most valuable victim / least valuable attacker) for captures; killer move heuristic for quiet moves
- Quiescence search: extend search at positions with pending captures until a quiet position is reached
- Time management: `context.WithDeadline` cancels search goroutine; main goroutine blocks on result channel
- UCI info emission: after each completed depth iteration, emit a formatted info line to the writer

---

## Alternatives Considered

### Alternative A: Fixed-Depth Minimax (No Iterative Deepening)

**Description**: Search to a fixed depth D, return the best move. No time management; depth is configured externally.

**Evaluation against requirements**:
- Time compliance: violates NFR-03 — fixed depth cannot guarantee time compliance; a complex position at depth 8 may take 10× longer than a simple one
- UCI: `go movetime <ms>` requires time-based termination, not depth-based
- Quality: does not scale to available time; weak on short time controls

**Rejection rationale**: Time compliance is a hard requirement (NFR-03). Fixed depth cannot satisfy it. Iterative deepening is the standard solution and is explicitly named in FR-02-01.

### Alternative B: Principal Variation Search (PVS / Negascout)

**Description**: Enhancement of alpha-beta that uses null-window searches for non-PV nodes, improving pruning efficiency. Used in most strong engines.

**Evaluation against requirements**:
- Performance: 20–30% more nodes searched per second vs plain alpha-beta (better pruning)
- Complexity: requires careful implementation of aspiration windows and re-search logic
- Correctness risk: subtle bugs in null-window logic are common

**Rejection rationale**: PVS is a v2 optimization. The 100k NPS target is achievable with plain alpha-beta at the 8x8 board level. PVS adds implementation complexity and correctness risk with no v1 requirement driving it. It is a natural upgrade once the basic alpha-beta is validated by perft and engine tests.

### Alternative C: Monte Carlo Tree Search (MCTS)

**Description**: Random playouts with UCB1 selection, used in AlphaZero-style engines. Does not require a hand-crafted evaluation function.

**Evaluation against requirements**:
- Complexity: requires neural network (Leela) or hand-tuned rollout policy; far exceeds v1 scope
- Evaluation: material + positional evaluation (FR-02-02, FR-02-03) implies alpha-beta; MCTS with hand-crafted rollouts is weaker than alpha-beta for standard material counting
- UCI info: MCTS produces different statistics (visits, confidence) that do not map naturally to UCI depth/score convention

**Rejection rationale**: MCTS is designed for domains where search-based evaluation is expensive (Go, complex games). For standard chess with a simple eval function, alpha-beta dominates MCTS. MCTS is an order of magnitude more complex to implement correctly.

---

## Consequences

**Positive**:
- Iterative deepening gives best-effort result at any time cutoff (always a move to return)
- Negamax frame halves the code size vs separate max/min functions
- Move ordering with MVV-LVA achieves near-optimal pruning for typical chess positions
- Quiescence search eliminates horizon effect on captures (required for tactically sound play)
- Context cancellation integrates cleanly with Go's concurrency model; no goroutine leak

**Negative**:
- Re-searching depths 1..N-1 wastes ~33% of nodes (standard ID cost, accepted by convention)
- Without a transposition table, the same positions are evaluated multiple times across iterations
- Killer moves are lost between moves (reset per search call); a global killer table would require shared mutable state (complexity vs benefit tradeoff accepted)

**v2 Upgrade Path**:
- Add transposition table (Zobrist hash → SearchResult cache): eliminates re-evaluation of identical positions across iterations
- Add aspiration windows: narrow alpha-beta window around expected score to reduce search tree
- Upgrade to PVS/Negascout: better pruning at non-PV nodes
- All upgrades are internal to `engine/search.go`; no public API changes

---

## Time Management Detail

The time manager computes an allocated time for the current move:

```
if MoveTime > 0:
    deadline = now + MoveTime

else:
    remaining = color == White ? WTime : BTime
    increment = color == White ? WInc : BInc
    allocated = remaining/30 + increment*0.8
    deadline = now + allocated
```

A `context.WithDeadline(ctx, deadline)` is created and passed to the search. The search checks `ctx.Done()` at the top of every node. On cancellation, it returns the best move from the last completed depth iteration. The check adds minimal overhead (channel receive is ~1ns).

The engine guarantees return within `deadline + 50ms` — the 50ms accounts for I/O overhead and OS scheduling jitter.
