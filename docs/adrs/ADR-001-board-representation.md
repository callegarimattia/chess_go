# ADR-001: Board Representation

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Morgan (solution architect)
**Affected components**: `internal/chess` (board.go, movegen.go)

---

## Context

The chess package requires a data structure to represent the placement of pieces on a 64-square board. The choice of board representation is the most consequential performance and complexity decision in the chess package: it determines how fast moves are generated, how much memory each position uses, and how complex the move generator is to implement and test.

Three families of board representation exist in chess programming:

1. **8x8 array** — A 64-element array indexed by square number. Piece lookup is O(1). Move generation iterates squares and applies directional offsets.
2. **Bitboards** — 12 separate 64-bit integers (one per piece type per color). Move generation uses bit manipulation (shifts, AND, OR, popcount). Achieves multi-million NPS. Standard in high-performance engines.
3. **0x88** — A 128-element array where the upper bit of the index indicates off-board. Simplifies boundary checking via bitwise AND. Obsolete in modern engines but was popular pre-2000.

The primary quality attributes driving this decision:
- **Correctness**: move generation must pass perft tests at depth 5 (AC-11)
- **Performance**: > 100,000 NPS on modern hardware (NFR-01)
- **Testability**: chess package >= 90% line coverage; move generator logic must be comprehensible
- **Maintainability**: solo/small team; complex bitboard magic requires deep expertise to debug

---

## Decision

Use an **8x8 array** (`[64]Piece`) as the board representation for v1.

The board is a value type: a `[64]byte` array embedded in `GameState`. Copying a position copies 64 bytes — trivially cheap. Square indexing: `a1=0, b1=1, ..., h8=63`. Piece encoding: one byte per square (zero = empty).

Move generation iterates piece positions and applies directional offsets. Boundary detection uses rank/file arithmetic. The legality filter applies each pseudo-legal move and checks whether the resulting position leaves the king in check.

---

## Alternatives Considered

### Alternative A: Bitboards

**Description**: 12 × 64-bit integers represent piece occupancy. Move generation uses pre-computed attack tables and bit manipulation. Standard in production engines (Stockfish, Leela).

**Evaluation against requirements**:
- Performance: achieves 1M–10M NPS; far exceeds the 100k NPS target
- Correctness: harder to implement; bugs are subtle (off-by-one in attack generation)
- Testability: complex bit manipulation is harder to unit test and reason about
- Maintainability: requires deep knowledge of magic bitboards or PEXT instructions to reach full performance

**Rejection rationale**: Bitboards are 10–100x more complex to implement correctly than an 8x8 array. The 100k NPS target is achievable with an 8x8 array (perft depth 5 = 4.8M nodes; at 100k NPS this takes ~48 seconds — acceptable for correctness testing). The performance headroom is not justified by v1 requirements. Bitboards are the correct v2 upgrade once profiling confirms the array is the bottleneck.

### Alternative B: 0x88 Representation

**Description**: 128-element array where square validity is tested by `sq & 0x88 == 0`. Move generation uses offset tables with inline validity checks.

**Evaluation against requirements**:
- Performance: comparable to 8x8 array; no significant advantage
- Correctness: boundary detection is simpler than 8x8 file arithmetic
- Maintainability: non-obvious indexing scheme; team unfamiliar with the representation

**Rejection rationale**: 0x88 provides minor boundary-detection simplification over 8x8 at the cost of a non-obvious 128-element layout. The 8x8 array with rank/file arithmetic is more readable and more commonly understood. No performance advantage justifies the unfamiliarity.

---

## Consequences

**Positive**:
- Simplest possible representation; move generator logic is straightforward to read and test
- No bit manipulation expertise required
- GameState copy is 64 bytes — fits in a cache line with header fields
- Easy to add debugging output (iterate squares and print piece names)

**Negative**:
- Move generation is slower than bitboards; generating all pseudo-legal moves requires iterating piece positions
- If future NPS requirements exceed 100k significantly, migration to bitboards will require rewriting movegen.go
- No pre-computed attack tables; each move candidate requires a direction loop

**Migration path to bitboards (v2)**:
- Bitboard representation lives entirely in `board.go` and `movegen.go`
- Public API (`GameState`, `Move`, `Game.LegalMoves()`, `Game.Apply()`) is unchanged
- Migration is internal to chess package with no consumer code changes

---

## Perft Performance Estimate (8x8 array)

Perft depth 5 = 4,865,609 nodes. At 100k NPS, completion time = ~49 seconds.
This is acceptable for a correctness test suite run during CI. It is not run on every `go test` invocation; it is gated behind a build tag or `-run Perft` flag.

The 100k NPS requirement (NFR-01) applies to engine search (which searches far fewer nodes per move decision), not perft.
