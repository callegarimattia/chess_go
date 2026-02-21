// Command chess-go is the TUI binary for the chess engine.
// It wires the tui, engine, and chess packages and starts the game loop.
package main

import (
	"os"
	"time"

	chess "chess_go/internal/chess"
	engine "chess_go/internal/engine"
	"chess_go/internal/tui"
)

func main() {
	game := tui.NewGame(os.Stdin, os.Stdout, func(g chess.Game, tc engine.TimeControl) chess.Move {
		tc = engine.TimeControl{MoveTime: 100 * time.Millisecond}
		return engine.Search(g, tc, os.Stderr).BestMove
	})
	game.Run()
}
