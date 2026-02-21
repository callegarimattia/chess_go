# Shared Artifacts Registry — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21

This registry tracks every data artifact that crosses component boundaries. Each entry has exactly one authoritative producer. All consumers must use the producer's output — no re-derivation.

---

## SA-01 · FEN String

| Field | Value |
|-------|-------|
| **ID** | SA-01 |
| **Name** | FEN string |
| **Description** | Forsyth-Edwards Notation — canonical position encoding |
| **Format** | `rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1` |
| **Producer** | `chess` package — `FENParser` |
| **Single source of truth** | `GameState.ToFEN()` method |
| **Consumers** | Board renderer (TUI), Board renderer (SSR), Engine search input, UCI position handler, PGN recorder |
| **Validation rules** | 6 space-separated fields; valid piece placement; active color `w`/`b`; castling rights subset of `KQkq` or `-`; en passant `-` or valid square; non-negative integers for clocks |
| **Error** | `ErrInvalidFEN` with field index and reason |

---

## SA-02 · Move Notation — UCI Format

| Field | Value |
|-------|-------|
| **ID** | SA-02 |
| **Name** | Move notation (UCI) |
| **Description** | Move in Universal Chess Interface format — from-square + to-square + optional promotion piece |
| **Format** | `e2e4`, `e7e8q`, `e1g1` (castling) |
| **Producer** | TUI input parser, SSR move endpoint, Engine search output |
| **Single source of truth** | `Move.UCIString()` method |
| **Consumers** | Move validator, Board state updater, PGN recorder (converts to SAN), UCI bestmove output |
| **Validation rules** | 4 or 5 characters; squares a1–h8; promotion piece one of `q r b n` |
| **Error** | `ErrInvalidMoveFormat` |

---

## SA-03 · Move Notation — SAN Format

| Field | Value |
|-------|-------|
| **ID** | SA-03 |
| **Name** | Move notation (SAN) |
| **Description** | Standard Algebraic Notation for human-readable display |
| **Format** | `e4`, `Nf3`, `O-O`, `Bxe5+`, `e8=Q#` |
| **Producer** | `chess` package — `move.ToSAN(gameState)` |
| **Single source of truth** | `Move.SANString(gs GameState)` method |
| **Consumers** | TUI move history display, SSR move history display, PGN recorder |
| **Validation rules** | Derived from legal move in context — never stored as input format |
| **Note** | SAN is output-only; all input uses UCI format (SA-02) |

---

## SA-04 · GameState

| Field | Value |
|-------|-------|
| **ID** | SA-04 |
| **Name** | GameState |
| **Description** | Complete, immutable snapshot of board position and game metadata |
| **Format** | Go struct |
| **Fields** | Board (piece placement), ActiveColor, CastlingRights, EnPassantSquare, HalfMoveClock, FullMoveNumber |
| **Producer** | `chess` package — `Game.Apply(move)` returns new GameState |
| **Single source of truth** | `Game` struct — immutable; each move produces a new GameState |
| **Consumers** | Move generator, Engine search, TUI renderer, SSR renderer, Result detector, PGN recorder |
| **Concurrency** | GameState is immutable — safe to share across goroutines |
| **Error** | Invalid transitions return error, never mutate existing state |

---

## SA-05 · PGN Record

| Field | Value |
|-------|-------|
| **ID** | SA-05 |
| **Name** | PGN record |
| **Description** | Portable Game Notation — complete game transcript |
| **Format** | Standard PGN text with headers and move list |
| **Producer** | `chess` package — `Game.ToPGN()` |
| **Single source of truth** | `Game.MoveHistory` + metadata headers |
| **Consumers** | TUI export, SSR export, external chess tools |
| **Required headers** | Event, Site, Date, Round, White, Black, Result |
| **Error** | `ErrPGNExport` if game is in invalid state |

---

## SA-06 · Engine Response

| Field | Value |
|-------|-------|
| **ID** | SA-06 |
| **Name** | Engine response |
| **Description** | Best move found by engine search, with optional analysis info |
| **Format (UCI)** | `bestmove e2e4 [ponder e7e5]` |
| **Format (info)** | `info depth 6 score cp 30 nodes 50000 nps 50000 pv g1f3 g8f6` |
| **Producer** | `engine` package — `Search(position, timeControl)` |
| **Single source of truth** | `SearchResult.BestMove` field |
| **Consumers** | UCI output handler, TUI game loop (applies engine move), SSR game loop |
| **Guarantee** | Always produces a legal move; never blocks past time limit + 50ms grace |
| **Error fallback** | If time expires before search starts: return first legal move |

---

## SA-07 · Time Control

| Field | Value |
|-------|-------|
| **ID** | SA-07 |
| **Name** | Time control |
| **Description** | Remaining time and increment for each side |
| **Format (UCI)** | `go wtime 300000 btime 300000 winc 0 binc 0` or `go movetime 1000` |
| **Format (internal)** | `TimeControl` struct with fields per side |
| **Producer** | TUI time manager, SSR time manager, UCI `go` command parser |
| **Single source of truth** | `TimeControl` struct passed to `Search()` |
| **Consumers** | Engine search (time manager), UCI handler |
| **Guarantee** | Engine must not exceed `movetime` or allocated time by more than 50ms |

---

## Artifact Flow Diagram

```
User Input (TUI/GUI/UCI)
         │
         ▼
    [SA-02 UCI Move]
         │
         ▼
   Move Validator ──── [SA-04 GameState] ◄──── FEN Parser ◄──── [SA-01 FEN]
         │                    │
         ▼                    ▼
   Board Updater        Move Generator
         │                    │
         ▼                    ▼
   [SA-04 GameState']   Engine Search ◄──── [SA-07 Time Control]
         │                    │
    ┌────┴────┐               ▼
    ▼         ▼         [SA-06 Engine Response]
TUI Render  SSR Render        │
    │         │               ▼
    ▼         ▼          Game Loop (applies bestmove)
[SA-03 SAN] [SA-03 SAN]      │
    │                         ▼
    └──────────────── [SA-05 PGN Record] ──► Export
```

---

## Cross-Reference: Artifacts by Component

| Component | Produces | Consumes |
|-----------|----------|----------|
| FEN parser | SA-04 | SA-01 |
| Move generator | SA-02 (list) | SA-04 |
| Board updater | SA-04 | SA-02, SA-04 |
| Engine search | SA-06 | SA-04, SA-07 |
| PGN recorder | SA-05 | SA-02, SA-03, SA-04 |
| TUI renderer | SA-03 | SA-04 |
| SSR renderer | SA-03 | SA-04 |
| UCI handler | SA-06 | SA-01, SA-02, SA-07 |
| TUI input | SA-02 | — |
| SSR endpoint | SA-02 | — |
