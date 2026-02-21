package tui

import (
	"io"
	"time"

	chess "chess_go/internal/chess"
	engine "chess_go/internal/engine"
)

// EngineFunc is the type of the engine callback injected into the TUI game loop.
type EngineFunc func(g chess.Game, tc engine.TimeControl) chess.Move

// Game manages the TUI game loop with injected I/O for testability.
type Game struct {
	r        io.Reader
	w        io.Writer
	engineFn EngineFunc
}

// NewGame constructs a TUI Game with the given I/O and engine function.
func NewGame(r io.Reader, w io.Writer, engineFn EngineFunc) Game {
	return Game{r: r, w: w, engineFn: engineFn}
}

// Run executes the walking skeleton game loop:
// 1. Load the starting position.
// 2. Call engineFn to select a move.
// 3. Apply the move.
// 4. Render the resulting board.
//
// For the walking skeleton, this executes one move and exits.
// Full interactive loop is implemented in US-24/US-25.
func (g Game) Run() {
	game, err := chess.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		return
	}

	tc := engine.TimeControl{MoveTime: 100 * time.Millisecond}
	move := g.engineFn(game, tc)

	newGame, err := game.Apply(move)
	if err != nil {
		return
	}

	Render(newGame, g.w)
}
