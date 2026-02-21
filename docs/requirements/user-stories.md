# User Stories — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21

Stories follow LeanUX format: `As a <persona>, I want <goal>, so that <outcome>`.
Each story is sized for a single iteration. Dependencies are noted.

---

## Feature 0 — Walking Skeleton

### US-00: End-to-end skeleton
**As** Daniel (engine developer),
**I want** the engine to select a random legal move from the starting position and display the updated board in the terminal,
**so that** I can validate that all architectural layers (chess package → engine → TUI) are connected correctly.

**Size**: S
**Depends on**: nothing
**Personas**: P1, P2, P3
**Artifacts produced**: SA-04 (GameState), SA-02 (Move), TUI render

---

## Epic 1 — Chess Logic Package

### US-01: Parse FEN
**As** Marco (library integrator),
**I want** to call `NewGameFromFEN(fen string)` and receive a valid `GameState`,
**so that** I can load any position into my application without reimplementing FEN parsing.

**Size**: S | **Depends on**: US-00

### US-02: Generate legal moves
**As** Marco (library integrator),
**I want** to call `game.LegalMoves()` and receive a complete list of legal moves,
**so that** I can build move selection logic without worrying about move legality.

**Size**: M | **Depends on**: US-01

### US-03: Apply a move
**As** Marco (library integrator),
**I want** to call `game.Apply(move)` and receive a new `GameState` without mutating the original,
**so that** I can safely explore positions in parallel and revert to previous states.

**Size**: S | **Depends on**: US-02

### US-04: Detect check and checkmate
**As** Marco (library integrator),
**I want** `game.InCheck()` and `game.Result()` to return accurate check/checkmate status,
**so that** my tournament manager can correctly end games and display results.

**Size**: M | **Depends on**: US-03

### US-05: Detect draw conditions
**As** Marco (library integrator),
**I want** `game.Result()` to detect fifty-move rule, threefold repetition, and insufficient material,
**so that** games end correctly without manual tracking in my application.

**Size**: M | **Depends on**: US-04

### US-06: Special moves — castling
**As** Marco (library integrator),
**I want** castling (e1g1, e1c1, e8g8, e8c8) to be generated and applied correctly,
**so that** positions involving castling behave identically to standard chess rules.

**Size**: M | **Depends on**: US-03

### US-07: Special moves — en passant
**As** Marco (library integrator),
**I want** en passant captures to be generated and applied correctly,
**so that** positions with en passant opportunity are handled without special-casing in my code.

**Size**: S | **Depends on**: US-03

### US-08: Special moves — pawn promotion
**As** Marco (library integrator),
**I want** pawn promotion moves (e7e8q, e7e8r, e7e8b, e7e8n) to be generated and applied correctly,
**so that** games can progress past promotion squares without errors.

**Size**: S | **Depends on**: US-03

### US-09: Export FEN and UCI notation
**As** Marco (library integrator),
**I want** `GameState.ToFEN()` and `Move.UCIString()` to produce standard notation,
**so that** I can interoperate with external chess tools and GUIs.

**Size**: S | **Depends on**: US-03

### US-10: Export SAN notation
**As** Marco (library integrator),
**I want** `Move.SANString(gs GameState)` to produce correct SAN notation,
**so that** I can display moves in a human-readable format in my tournament manager.

**Size**: S | **Depends on**: US-09

### US-11: Export PGN
**As** Marco (library integrator),
**I want** `Game.ToPGN()` to produce a valid PGN record,
**so that** I can archive game transcripts and replay them in standard chess software.

**Size**: S | **Depends on**: US-10

### US-12: Perft validation
**As** Daniel (engine developer),
**I want** the move generator to pass perft tests up to depth 5 against known values,
**so that** I can be confident the move generation is bug-free before tuning the engine.

**Size**: M | **Depends on**: US-06, US-07, US-08

---

## Epic 2 — Engine

### US-13: Random move engine (skeleton)
**As** Daniel (engine developer),
**I want** a minimal engine that selects a random legal move,
**so that** the walking skeleton (US-00) can demonstrate the full pipeline without requiring search logic.

**Size**: XS | **Depends on**: US-02
**Note**: Part of US-00 walking skeleton; extracted as its own story for clarity.

### US-14: Alpha-beta search
**As** Daniel (engine developer),
**I want** the engine to implement alpha-beta search with iterative deepening,
**so that** it finds stronger moves as search time increases.

**Size**: L | **Depends on**: US-12

### US-15: Material evaluation
**As** Daniel (engine developer),
**I want** the engine to evaluate positions using standard piece values (P=100, N=320, B=330, R=500, Q=900),
**so that** it makes materially sound decisions.

**Size**: S | **Depends on**: US-14

### US-16: Positional evaluation
**As** Daniel (engine developer),
**I want** the engine to use piece-square tables for positional bonuses,
**so that** it develops pieces to strong squares and avoids passive positions.

**Size**: M | **Depends on**: US-15

### US-17: Time management
**As** Daniel (engine developer),
**I want** the engine to respect `movetime` and `wtime/btime` UCI time controls,
**so that** it never loses on time and plays within allocated limits.

**Size**: M | **Depends on**: US-14

### US-18: Quiescence search
**As** Daniel (engine developer),
**I want** the engine to extend search at tactical positions (captures, checks),
**so that** it avoids the horizon effect and evaluates quiet positions correctly.

**Size**: M | **Depends on**: US-14

### US-19: Move ordering
**As** Daniel (engine developer),
**I want** the engine to order captures before quiet moves and try killer moves early,
**so that** alpha-beta pruning is more effective and the engine searches deeper.

**Size**: M | **Depends on**: US-14

---

## Epic 3 — UCI Protocol

### US-20: UCI handshake
**As** Daniel (engine developer),
**I want** the binary to respond correctly to `uci`, `isready`, and `ucinewgame`,
**so that** I can connect it to standard chess GUIs like Arena or Cute Chess.

**Size**: S | **Depends on**: US-13

### US-21: Position command
**As** Daniel (engine developer),
**I want** the engine to handle `position startpos [moves ...]` and `position fen <fen> [moves ...]`,
**so that** I can set up any position before searching.

**Size**: S | **Depends on**: US-20, US-09

### US-22: Go command
**As** Daniel (engine developer),
**I want** the engine to handle `go movetime <ms>` and `go wtime <ms> btime <ms>`,
**so that** it searches within the time budget I provide.

**Size**: S | **Depends on**: US-21, US-17

### US-23: Stop and quit commands
**As** Daniel (engine developer),
**I want** the engine to handle `stop` (emit bestmove immediately) and `quit` (exit cleanly),
**so that** I can control the engine lifecycle from a GUI or test harness.

**Size**: S | **Depends on**: US-22

---

## Epic 4 — TUI

### US-24: Render board in terminal
**As** Sofia (TUI player),
**I want** to see an ASCII chess board with piece symbols and coordinates when I launch `chess-go`,
**so that** I can understand the position immediately without reading documentation.

**Size**: S | **Depends on**: US-00

### US-25: Accept and validate move input
**As** Sofia (TUI player),
**I want** to type a move in UCI format and see it applied or receive a clear error,
**so that** I can play moves confidently without worrying about input format details.

**Size**: S | **Depends on**: US-24, US-03

### US-26: Engine response in TUI
**As** Sofia (TUI player),
**I want** to see the engine's response move applied to the board after I move,
**so that** I can play a complete game against the engine.

**Size**: S | **Depends on**: US-25, US-22

### US-27: Game status in TUI
**As** Sofia (TUI player),
**I want** to see check notifications, game results, and draw reasons displayed clearly,
**so that** I always know the state of the game without counting moves manually.

**Size**: S | **Depends on**: US-26, US-04, US-05

### US-28: PGN export from TUI
**As** Sofia (TUI player),
**I want** to save the game as a PGN file when the game ends,
**so that** I can review my games later in a standard chess tool.

**Size**: S | **Depends on**: US-27, US-11

---

## Epic 5 — SSR GUI

### US-29: Serve HTML board
**As** Priya (web player),
**I want** to open a browser and see a chess board at `http://localhost:8080`,
**so that** I can play chess without installing any software beyond the binary.

**Size**: S | **Depends on**: US-24 (board rendering logic)

### US-30: Create game session
**As** Priya (web player),
**I want** to click "New Game" and see a fresh board in the browser,
**so that** I can start playing immediately.

**Size**: S | **Depends on**: US-29

### US-31: Make a move via SSR
**As** Priya (web player),
**I want** to click a piece and a destination square to make a move,
**so that** I can play chess without typing notation.

**Size**: M | **Depends on**: US-30, US-25

### US-32: Engine response in SSR
**As** Priya (web player),
**I want** the engine to respond immediately after my move and the updated board to render,
**so that** the game flows naturally without manual refresh.

**Size**: S | **Depends on**: US-31, US-26

### US-33: Game result in SSR
**As** Priya (web player),
**I want** to see a clear result message and a "New Game" link when the game ends,
**so that** I can play again without refreshing or navigating manually.

**Size**: S | **Depends on**: US-32, US-27

---

## Story Map Summary

```
                  SKELETON    CHESS PKG    ENGINE    UCI     TUI     SSR
Walking skeleton    US-00
FEN/move/state               US-01..03
Special moves                US-06..08
Notation                     US-09..11
Perft                        US-12
Search                                    US-14..16
Time mgmt                                US-17
Quiescence                               US-18..19
UCI protocol                                       US-20..23
TUI game loop                                              US-24..28
SSR GUI                                                            US-29..33
```
