# Requirements — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21 | **Status**: Draft

---

## Project Vision

Build a chess engine in Go comprising three layers:
1. **`chess` package** — reusable Go library for chess logic (rules, move generation, game state)
2. **`engine` package** — standalone search engine (alpha-beta, evaluation, time management)
3. **Interfaces** — TUI (terminal) and SSR GUI (web), both thin layers over the packages above

The system is architecturally layered: the `chess` package has no dependency on the engine; the engine depends only on the `chess` package; TUI and SSR depend on both.

---

## Scope

### In Scope
- Complete chess rules implementation (all legal moves, special moves, draw conditions)
- Alpha-beta search with iterative deepening
- Material + positional evaluation function
- UCI protocol compliance
- TUI: ASCII board, move input, game loop, PGN export
- SSR GUI: HTTP server, HTML board, move via form POST, game loop
- Perft test suite for move generation validation

### Out of Scope (v1)
- Opening books
- Endgame tablebases
- Multi-PV analysis
- Online play / multiplayer
- Time control UI (TUI accepts movetime only in v1)
- Authentication or sessions (SSR is local-only in v1)

---

## Functional Requirements

### FR-01: Chess Logic Package

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01-01 | Parse FEN strings into GameState | Must |
| FR-01-02 | Generate all legal moves from any position | Must |
| FR-01-03 | Apply a move to produce a new GameState (immutable) | Must |
| FR-01-04 | Detect check, checkmate, and stalemate | Must |
| FR-01-05 | Detect draw: fifty-move rule, threefold repetition, insufficient material | Must |
| FR-01-06 | Support special moves: castling (K/Q-side), en passant, pawn promotion | Must |
| FR-01-07 | Export GameState as FEN string | Must |
| FR-01-08 | Export move as UCI string | Must |
| FR-01-09 | Export move as SAN string given context GameState | Must |
| FR-01-10 | Export game as PGN | Should |
| FR-01-11 | Validate FEN string and return typed errors | Must |

### FR-02: Engine

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-02-01 | Implement alpha-beta search with iterative deepening | Must |
| FR-02-02 | Implement basic material evaluation (piece values) | Must |
| FR-02-03 | Implement positional evaluation (piece-square tables) | Should |
| FR-02-04 | Respect time controls: movetime, wtime/btime with increment | Must |
| FR-02-05 | Always return bestmove within time limit + 50ms grace | Must |
| FR-02-06 | Emit UCI info lines (depth, score, nodes, nps, pv) during search | Must |
| FR-02-07 | Support quiescence search to avoid horizon effect | Should |
| FR-02-08 | Support move ordering (captures first, killer moves) | Should |

### FR-03: UCI Protocol

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-03-01 | Respond to `uci` with `id name`, `id author`, `uciok` | Must |
| FR-03-02 | Respond to `isready` with `readyok` | Must |
| FR-03-03 | Handle `ucinewgame` by resetting internal state | Must |
| FR-03-04 | Handle `position startpos [moves ...]` | Must |
| FR-03-05 | Handle `position fen <fen> [moves ...]` | Must |
| FR-03-06 | Handle `go movetime <ms>` | Must |
| FR-03-07 | Handle `go wtime <ms> btime <ms> [winc <ms> binc <ms>]` | Must |
| FR-03-08 | Handle `stop` and emit bestmove within 100ms | Must |
| FR-03-09 | Handle `quit` by exiting cleanly | Must |
| FR-03-10 | Ignore unknown commands gracefully (no crash) | Must |

### FR-04: TUI

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-04-01 | Render ASCII board with pieces and coordinates | Must |
| FR-04-02 | Accept move input in UCI format (e.g. e2e4) | Must |
| FR-04-03 | Display "White/Black to move" indicator | Must |
| FR-04-04 | Display "Engine thinking..." while engine searches | Must |
| FR-04-05 | Display illegal move message and re-prompt | Must |
| FR-04-06 | Display check notification after move | Must |
| FR-04-07 | Display game result (checkmate, stalemate, draw + reason) | Must |
| FR-04-08 | Prompt to save PGN at game end | Should |
| FR-04-09 | Display move history (scrollable) | Should |
| FR-04-10 | Support configurable engine search depth or movetime | Could |

### FR-05: SSR GUI

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-05-01 | Serve HTML chess board at GET / | Must |
| FR-05-02 | Create new game session via POST /game/new | Must |
| FR-05-03 | Render full board page at GET /game/{id} | Must |
| FR-05-04 | Accept move via POST /game/{id}/move?from=e2&to=e4 | Must |
| FR-05-05 | Validate move server-side; return HTTP 422 on illegal move | Must |
| FR-05-06 | Engine responds immediately after player move (same request cycle) | Must |
| FR-05-07 | Display check notification in rendered board | Must |
| FR-05-08 | Display game result and offer new game link | Must |
| FR-05-09 | Game session persists in memory (no external DB in v1) | Must |
| FR-05-10 | Board rendered as HTML table or SVG (no JS required) | Must |

---

## Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-01 | Engine search speed | > 100,000 nodes/second on modern hardware |
| NFR-02 | Move generation correctness | Perft results match known values at depth 5 |
| NFR-03 | Time control compliance | bestmove always within movetime + 50ms |
| NFR-04 | Package API stability | No breaking changes within v0.x minor versions |
| NFR-05 | Concurrency | GameState is immutable and goroutine-safe |
| NFR-06 | Zero dependencies (chess package) | chess package has no external dependencies |
| NFR-07 | SSR response time | Board page renders in < 200ms including engine move |
| NFR-08 | TUI render latency | Board redraws in < 50ms |
| NFR-09 | Test coverage | chess package >= 90% line coverage |
| NFR-10 | Build | `go build ./...` succeeds on clean checkout with `go 1.22+` |

---

## Constraints

- **Language**: Go 1.22+
- **No external dependencies** in the `chess` package
- Engine and interface layers may use standard library and minimal well-known dependencies
- SSR GUI: no JavaScript required for core functionality
- All interfaces must work on Linux, macOS, and Windows

---

## Architecture Overview

```
cmd/
  chess-go/       ← TUI binary entry point
  chess-server/   ← SSR GUI binary entry point

internal/
  chess/          ← chess logic package (zero deps)
    board.go
    move.go
    game.go
    fen.go
    movegen.go
    result.go
    pgn.go

  engine/         ← search engine (depends on chess/)
    search.go
    eval.go
    time.go
    uci.go

  tui/            ← terminal UI (depends on chess/, engine/)
    renderer.go
    input.go
    loop.go

  web/            ← SSR HTTP handlers (depends on chess/, engine/)
    handler.go
    session.go
    template.go
```
