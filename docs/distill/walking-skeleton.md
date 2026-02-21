# Walking Skeleton Implementation Guide
**Epic**: chess-engine | **Feature 0** | **Date**: 2026-02-21

This document guides the software-crafter through implementing Feature 0: the walking skeleton that validates all architectural layers are connected end-to-end.

---

## Goal

Deliver the minimum slice of production code that makes `TestWalkingSkeleton_PipelineConnectsAllLayers` pass with real implementation (not just compilation).

Observable outcome: `chess-go` binary launched in a terminal selects a random legal move from the starting position, displays the updated ASCII board, and exits with code 0.

User story: US-00.

---

## Pipeline (from architecture-design.md §13)

```
Input:    starting FEN (hardcoded)
chess:    NewGameFromFEN → GameState
chess:    GameState.LegalMoves() → []Move (20 moves)
engine:   random selection from []Move → Move (no search logic)
chess:    GameState.Apply(Move) → newGameState
tui:      Renderer.Render(newGameState) → ASCII board to stdout
Exit:     code 0
```

---

## Implementation Steps

### Step 1: Create go.mod

```
module chess_go
go 1.22
```

### Step 2: Implement internal/chess (minimum surface for skeleton)

Files to create:

- `internal/chess/board.go` — `type Board [64]Piece`, `type Square uint8`, `type Piece uint8` constants, `type Color uint8`
- `internal/chess/move.go` — `type Move struct { From, To Square; Promotion Piece }`, `UCIString() string`
- `internal/chess/game.go` — `type GameState struct { ... }`, `type Game struct { State GameState; ... }`, `LegalMoves() []Move`, `Apply(Move) (Game, error)`, `InCheck() bool`, `Result() GameResult`, `ToFEN() string`
- `internal/chess/fen.go` — `NewGameFromFEN(fen string) (Game, error)`, `ErrInvalidFEN`
- `internal/chess/movegen.go` — generates all pseudo-legal moves, filters for legality

Minimum correctness for the skeleton:
- `NewGameFromFEN(StartingFEN)` must return a valid Game
- `LegalMoves()` from the starting position must return exactly 20 moves
- `Apply(move)` must update the board correctly

### Step 3: Implement internal/engine (minimum: random move)

File: `internal/engine/search.go`

```go
// Search selects a random legal move (skeleton implementation).
// Replace with alpha-beta in US-14.
func Search(g chess.Game, tc TimeControl, info io.Writer) SearchResult {
    moves := g.LegalMoves()
    if len(moves) == 0 {
        return SearchResult{}
    }
    return SearchResult{BestMove: moves[rand.Intn(len(moves))]}
}
```

No time management, no evaluation, no info lines for the skeleton.

### Step 4: Implement internal/tui (minimum: render board)

File: `internal/tui/renderer.go`

```go
// Render writes an ASCII board representation of g to w.
func Render(g chess.Game, w io.Writer) {
    // 8x8 grid with piece symbols and coordinates.
    // Minimum: piece letters (K Q R B N P for White, k q r b n p for Black),
    //          rank numbers 1-8 on the left, file letters a-h at the bottom.
}
```

File: `internal/tui/loop.go`

```go
type Game struct { r io.Reader; w io.Writer; engineFn EngineFunc }
func NewGame(r io.Reader, w io.Writer, engineFn EngineFunc) Game { ... }
func (g Game) Run() { ... }
```

For the skeleton, Run() can:
1. Load the starting position
2. Call engineFn to get a move
3. Apply the move
4. Render the board
5. Return

### Step 5: Implement cmd/chess-go

File: `cmd/chess-go/main.go`

```go
func main() {
    game := tui.NewGame(os.Stdin, os.Stdout, func(g chess.Game, tc engine.TimeControl) chess.Move {
        return engine.Search(g, tc, os.Stderr).BestMove
    })
    game.Run()
}
```

### Step 6: Build and verify

```bash
go build ./cmd/chess-go
./chess-go
```

Expected output: an ASCII board after one random legal move, then exit 0.

---

## Acceptance Gate

Remove `t.Skip()` from `TestWalkingSkeleton_PipelineConnectsAllLayers` and run:

```bash
go test ./tests/acceptance/chess-engine/steps/... -run TestWalkingSkeleton -v
```

The test should PASS with real production code driving all pipeline stages.

---

## Skeleton Litmus Test

Before declaring the skeleton done, answer: "Can Daniel (engine developer) launch the binary, see a legal move played, and confirm the board is updated?"

If yes: the skeleton delivers observable user value end-to-end and the outer TDD loop can begin.

---

## Dependencies Enabled by the Skeleton

Once the skeleton passes, these stories become unblocked:

- US-01 (FEN parsing) — `NewGameFromFEN` exists
- US-02 (move generation) — `LegalMoves()` exists
- US-03 (apply move) — `Apply()` exists
- US-24 (TUI rendering) — `tui.Render()` exists
- US-13 (random engine) — `engine.Search()` skeleton exists

The next scenario to enable after the skeleton: `TestFENParser_ValidStartingPosition`.
