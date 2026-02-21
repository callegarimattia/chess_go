# Definition of Ready Checklist — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21 | **Validated by**: Luna (nw-product-owner)

The DoR gate must pass before handing off to the DESIGN wave (nw-solution-architect).
All 8 items must have evidence. Items marked FAIL block the handoff.

---

## DoR Item 1: Problem Statement is Clear

**Criterion**: The feature's purpose and value are unambiguous. We know what problem we're solving and for whom.

**Evidence**:
- 4 personas defined with specific goals: Marco (library integration), Sofia (TUI play), Daniel (engine development), Priya (web play)
- Journey maps document each persona's entry point, goal, and emotional arc
- Problem: no idiomatic Go chess engine exists that cleanly separates reusable library from interface layers

**Status**: PASS

---

## DoR Item 2: User Stories Follow Invest Criteria

**Criterion**: Each story is Independent, Negotiable, Valuable, Estimable, Small, Testable.

**Evidence** (sampled):

| Story | I | N | V | E | S | T |
|-------|---|---|---|---|---|---|
| US-00 Walking Skeleton | ✓ | ✓ | ✓ (validates arch) | ✓ (S) | ✓ | ✓ (AC-00) |
| US-01 FEN Parser | ✓ | ✓ | ✓ (all layers need it) | ✓ (S) | ✓ | ✓ (AC-01) |
| US-14 Alpha-beta | ✓ | ✓ | ✓ (engine quality) | ✓ (L) | ✓ (1 sprint) | ✓ (AC-12) |
| US-31 SSR move | ✓ | ✓ | ✓ (web player UX) | ✓ (M) | ✓ | ✓ (AC-16) |

All 34 stories reviewed. No story bundles multiple unrelated concerns.

**Status**: PASS

---

## DoR Item 3: Acceptance Criteria are Testable

**Criterion**: Every acceptance criterion can be verified by an automated test or a deterministic manual procedure.

**Evidence**:
- 16 acceptance criterion groups (AC-00 through AC-16)
- Each criterion specifies Given/When/Then with concrete values (e.g. "exactly 20 moves", "within 1050ms", "HTTP 422")
- Perft values in AC-11 are internationally known reference values
- No subjective criteria ("feels fast", "looks nice") — all quantified

**Status**: PASS

---

## DoR Item 4: Dependencies are Identified and Resolved

**Criterion**: External dependencies (tools, data, teams, decisions) are known and unblocked.

**Evidence**:
- All 34 stories have explicit dependency references (e.g. US-14 depends on US-12)
- Dependency graph is a DAG with US-00 and US-01 at the root — no circular deps
- External dependencies: Go 1.22+ (available), no external chess libraries required
- Architecture decisions deferred to DESIGN wave (board representation, search algorithm details) — correctly scoped out of requirements
- No team or infrastructure dependencies (solo/small team project)

**Status**: PASS

---

## DoR Item 5: UX Journey is Complete

**Criterion**: Happy path, error paths, and emotional arc are documented for all affected personas.

**Evidence**:
- `docs/ux/chess-engine/journey-chess-engine-visual.md`: visual journeys for all 4 personas
- `docs/ux/chess-engine/journey-chess-engine.yaml`: YAML schema with emotional arcs and scores
- Emotional arcs validated: all personas start at score ≤ 3 (skeptical/curious) and end at ≥ 4 (satisfied/trusting)
- Error paths: 6 documented in YAML (EP-01 through EP-06) covering illegal moves, invalid FEN, engine timeout, missing result detection, castling through check, and draw conditions
- Screen mockups included for TUI (launch, after-move) and UCI session transcript

**Status**: PASS

---

## DoR Item 6: Shared Artifacts are Tracked

**Criterion**: Every data artifact that crosses a component boundary has exactly one producer and is registered.

**Evidence**:
- `docs/ux/chess-engine/shared-artifacts-registry.md`: 7 artifacts registered (SA-01 through SA-07)
- Each artifact has: ID, name, description, format, producer, single source of truth, consumers, validation rules
- Artifact flow diagram confirms no artifact has two producers
- Potential conflict resolved: UCI input uses SA-02 (UCI notation); display uses SA-03 (SAN) — clearly separated

**Status**: PASS

---

## DoR Item 7: Stories are Appropriately Sized

**Criterion**: No story is larger than one sprint. L stories are the maximum allowed (decomposed further if needed).

**Evidence**:
- Story sizes: XS(1), S(22), M(9), L(2), XL(0)
- US-14 (alpha-beta) is L — acceptable; alpha-beta is a single well-understood algorithm
- US-16 (positional eval) is M — acceptable; piece-square tables are mechanical
- No XL stories
- L stories have well-defined completion criteria in acceptance-criteria.md

**Status**: PASS

---

## DoR Item 8: Handoff Package is Complete

**Criterion**: All artifacts required by the DESIGN wave are present and cross-referenced.

**Checklist**:

| Artifact | Path | Status |
|----------|------|--------|
| Visual journey map | `docs/ux/chess-engine/journey-chess-engine-visual.md` | ✓ Created |
| Journey YAML schema | `docs/ux/chess-engine/journey-chess-engine.yaml` | ✓ Created |
| Gherkin scenarios | `docs/ux/chess-engine/journey-chess-engine.feature` | ✓ Created |
| Shared artifacts registry | `docs/ux/chess-engine/shared-artifacts-registry.md` | ✓ Created |
| Requirements | `docs/requirements/requirements.md` | ✓ Created |
| User stories | `docs/requirements/user-stories.md` | ✓ Created |
| Acceptance criteria | `docs/requirements/acceptance-criteria.md` | ✓ Created |
| DoR checklist | `docs/requirements/dor-checklist.md` | ✓ This file |

All 8 artifacts present. No broken cross-references detected.

**Status**: PASS

---

## DoR Gate Result

| Item | Status |
|------|--------|
| 1. Problem statement clear | PASS |
| 2. Stories follow INVEST | PASS |
| 3. Acceptance criteria testable | PASS |
| 4. Dependencies identified | PASS |
| 5. UX journey complete | PASS |
| 6. Shared artifacts tracked | PASS |
| 7. Stories appropriately sized | PASS |
| 8. Handoff package complete | PASS |

**GATE: PASS — Ready for DESIGN wave handoff to nw-solution-architect**

---

## Handoff Summary for Solution Architect

### What to design:

1. **Board representation**: bitboard vs 8x8 array vs 0x88 — tradeoff: speed (bitboard) vs simplicity (8x8). Recommendation: 8x8 for v1, bitboard for v2.
2. **Search algorithm**: alpha-beta with iterative deepening (US-14) — standard, well-understood.
3. **Package layout**: `internal/chess`, `internal/engine`, `internal/tui`, `internal/web` as proposed in requirements.md.
4. **SSR rendering**: Go `html/template` + HTTP server, no JS. Session in-memory map.
5. **TUI**: `golang.org/x/term` or raw `os.Stdin` reads — no external TUI framework required.
6. **UCI binary**: reads from `os.Stdin`, writes to `os.Stdout`, engine in goroutine.

### Key constraints for architect:
- `chess` package: zero external dependencies
- GameState must be immutable (new state per move)
- Engine must always return within `movetime + 50ms`
- Go 1.22+ module
