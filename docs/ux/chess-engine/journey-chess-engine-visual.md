# UX Journey Map — Chess Engine in Go
**Epic**: chess-engine | **Research depth**: deep-dive | **Date**: 2026-02-21

---

## Personas

| ID | Name | Role | Goal |
|----|------|------|------|
| P1 | Marco Rossi | Library integrator | Import `chess_go` package, call clean API, embed in tournament manager |
| P2 | Sofia Chen | Human TUI player | Open terminal, play a casual game against the engine |
| P3 | Daniel Okafor | Engine developer | Tune search/eval, benchmark NPS, run UCI protocol headless |
| P4 | Priya Nair | Casual web player | Open browser, play against engine via SSR GUI |

---

## Journey 1: Marco Rossi — Library Integrator

```
STAGE        | Discover          | Setup             | First Move        | Integrate         | Trust
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ACTION       | Finds pkg on      | go get chess_go   | Calls             | Wires into        | Ships tournament
             | pkg.go.dev        | reads godoc       | game.Move("e2e4") | game loop         | manager
-------------|-------------------|-------------------|-------------------|-------------------|------------------
TOUCHPOINT   | godoc / README    | terminal + editor | REPL / test file  | own codebase      | production
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ARTIFACT     |                   | go.mod updated    | Move struct       | GameState         | PGN export
             |                   |                   | returned          | embedded          |
-------------|-------------------|-------------------|-------------------|-------------------|------------------
EMOTION      | Skeptical         | Cautious          | Relieved          | Confident         | Trusting
             | "edge cases?"     | "docs complete?"  | "it just works"   | "API is clean"    | "shipping this"
-------------|-------------------|-------------------|-------------------|-------------------|------------------
RISK         | Unclear API       | Missing examples  | Panic on illegal  | Global state leak | Silent bugs
```

**Emotional arc**: ↓ Skeptical → ↑ Exploratory → ↑↑ Confident → ↑↑↑ Trusting

**Error paths**:
- Illegal move: must return typed `ErrIllegalMove`, never panic
- Bad FEN input: `ErrInvalidFEN` with position of parse failure
- Nil game state: library must be safe for concurrent use (or document it is not)

---

## Journey 2: Sofia Chen — Human TUI Player

```
STAGE        | Launch            | Orientation       | First Move        | Mid-Game          | End
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ACTION       | $ chess-go        | Reads board       | Types "e2e4"      | Plays 20 moves    | Checkmate shown
             |                   | ASCII render      | Engine responds   | Sees eval hint    | Game saved
-------------|-------------------|-------------------|-------------------|-------------------|------------------
TOUCHPOINT   | Terminal          | TUI board render  | Move prompt       | Turn loop         | Result + PGN
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ARTIFACT     |                   | Board display     | Updated board     | Move history      | PGN file
-------------|-------------------|-------------------|-------------------|-------------------|------------------
EMOTION      | Curious           | Oriented          | Engaged           | Tense / Focused   | Satisfied
             | "will it work?"   | "I see the board" | "it responded!"   | "this is real"    | "good game"
-------------|-------------------|-------------------|-------------------|-------------------|------------------
RISK         | Garbled render    | Wrong coordinates | Illegal move msg  | Engine hangs      | No result shown
```

**TUI Screen Mockup — Initial Board**:
```
  chess-go v0.1  [White to move]

  8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜
  7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟
  6 · · · · · · · ·
  5 · · · · · · · ·
  4 · · · · · · · ·
  3 · · · · · · · ·
  2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙
  1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖
    a b c d e f g h

  Your move (e.g. e2e4): _
```

**TUI Screen Mockup — After Move**:
```
  chess-go v0.1  [Black to move]  Engine thinking...

  8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜
  7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟
  6 · · · · · · · ·
  5 · · · · · · · ·
  4 · · · · ♙ · · ·
  3 · · · · · · · ·
  2 ♙ ♙ ♙ ♙ · ♙ ♙ ♙
  1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖
    a b c d e f g h

  White: e2e4  |  Moves: 1
```

**Emotional arc**: → Curious → ↑ Oriented → ↑↑ Engaged → ↑↑ Tense → ↑↑↑ Satisfied

**Error paths**:
- Illegal move: `Illegal move: e2e5. Try again: _`
- Engine timeout: fallback to random legal move, display `[Engine: time limit reached]`
- Checkmate: display `Checkmate! Black wins.` then prompt to save PGN

---

## Journey 3: Daniel Okafor — Engine Developer

```
STAGE        | Bootstrap         | UCI Handshake     | Position Test     | Benchmark         | Tune
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ACTION       | Build binary      | Send uci/uciok    | position fen +    | go perft 5        | Adjust eval
             | $ go build        | isready/readyok   | go movetime 1000  | measure NPS       | weights
-------------|-------------------|-------------------|-------------------|-------------------|------------------
TOUCHPOINT   | Terminal / IDE    | stdin/stdout pipe | Chess GUI (Arena) | Benchmark harness | Source code
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ARTIFACT     | chess-go binary   | UCI session       | bestmove e2e4     | NPS report        | tuned binary
-------------|-------------------|-------------------|-------------------|-------------------|------------------
EMOTION      | Pragmatic         | Validating        | Testing           | Iterating         | Satisfied
             | "does it build?"  | "UCI compliant?"  | "correct move?"   | "fast enough?"    | "tournament ready"
-------------|-------------------|-------------------|-------------------|-------------------|------------------
RISK         | Build fails       | UCI non-compliant | Wrong bestmove    | Too slow          | Regression
```

**UCI session mockup**:
```
< uci
> id name chess-go
> id author chess-go contributors
> uciok
< isready
> readyok
< position startpos moves e2e4 e7e5
< go movetime 1000
> info depth 6 score cp 30 nodes 50000 nps 50000 pv g1f3
> bestmove g1f3
```

**Emotional arc**: → Pragmatic → ↑ Validating → ↑↑ Iterating → ↑↑↑ Satisfied

**Error paths**:
- UCI non-compliance: engine must respond to all required UCI commands
- Bestmove timeout: must always emit `bestmove` within time limit + 50ms grace
- Invalid position FEN: log to stderr, emit `bestmove 0000` (null move) as safe fallback

---

## Journey 4: Priya Nair — Casual Web Player (SSR GUI)

```
STAGE        | Visit             | Start Game        | Make Move         | Mid-Game          | End
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ACTION       | Opens browser     | Clicks "New Game" | Clicks piece      | Plays turns       | Sees result
             | localhost:8080    |                   | clicks target sq  |                   | offered rematch
-------------|-------------------|-------------------|-------------------|-------------------|------------------
TOUCHPOINT   | Browser / HTML    | SSR page reload   | POST /move        | Page renders      | Result page
-------------|-------------------|-------------------|-------------------|-------------------|------------------
ARTIFACT     |                   | Game session ID   | Updated board     | Move history      | Final position
-------------|-------------------|-------------------|-------------------|-------------------|------------------
EMOTION      | Casual            | Ready             | Excited           | Engaged           | Complete
             | "just want play"  | "board is clear"  | "piece moved!"    | "good game"       | "rematch?"
-------------|-------------------|-------------------|-------------------|-------------------|------------------
RISK         | Slow SSR          | Session lost      | Move rejected     | Engine lag        | No rematch option
```

**SSR GUI — board render (HTML)**:
```html
<!-- Server renders complete page on each move -->
GET  /           → start page
POST /game/new   → creates session, redirects to /game/{id}
GET  /game/{id}  → renders full board HTML
POST /game/{id}/move?from=e2&to=e4 → validates, engine responds, redirects
GET  /game/{id}  → renders updated board
```

**Emotional arc**: → Casual → ↑ Ready → ↑↑ Excited → ↑↑ Engaged → ↑ Complete

---

## Cross-Persona Shared Artifacts Flow

```
FEN string ──────────────────────────────────────────────┐
                                                          ↓
Input (TUI/GUI/UCI) → Parser → Board State → Move Generator → Engine Search
                                    ↓                              ↓
                              Display (TUI/GUI)            bestmove (UCI/API)
                                    ↓
                              PGN recorder → export
```
