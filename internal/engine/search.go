// Package engine implements chess search and time management.
// For the walking skeleton, Search selects a random legal move.
// Alpha-beta with iterative deepening is implemented in US-14.
package engine

import (
	"io"
	"math/rand"
	"time"

	chess "chess_go/internal/chess"
)

// TimeControl specifies how long the engine may think.
type TimeControl struct {
	MoveTime time.Duration // exact time for this move; 0 = use wtime/btime
	WTime    time.Duration // White remaining time
	BTime    time.Duration // Black remaining time
	WInc     time.Duration // White increment per move
	BInc     time.Duration // Black increment per move
}

// SearchResult holds the result of a search.
type SearchResult struct {
	BestMove chess.Move
	Score    int
	Depth    int
	Nodes    int64
	Elapsed  time.Duration
}

// Search selects a random legal move (skeleton implementation).
// Replace with alpha-beta in US-14.
func Search(g chess.Game, tc TimeControl, info io.Writer) SearchResult {
	start := time.Now()
	moves := g.LegalMoves()
	if len(moves) == 0 {
		return SearchResult{Elapsed: time.Since(start)}
	}
	//nolint:gosec
	best := moves[rand.Intn(len(moves))]
	return SearchResult{
		BestMove: best,
		Depth:    1,
		Nodes:    int64(len(moves)),
		Elapsed:  time.Since(start),
	}
}
