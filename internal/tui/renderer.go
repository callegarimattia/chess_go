// Package tui implements the terminal game loop and ASCII board renderer.
package tui

import (
	"fmt"
	"io"

	chess "chess_go/internal/chess"
)

// pieceASCII maps a Piece to its ASCII character for display.
func pieceASCII(p chess.Piece) byte {
	switch p {
	case chess.WhitePawn:
		return 'P'
	case chess.WhiteKnight:
		return 'N'
	case chess.WhiteBishop:
		return 'B'
	case chess.WhiteRook:
		return 'R'
	case chess.WhiteQueen:
		return 'Q'
	case chess.WhiteKing:
		return 'K'
	case chess.BlackPawn:
		return 'p'
	case chess.BlackKnight:
		return 'n'
	case chess.BlackBishop:
		return 'b'
	case chess.BlackRook:
		return 'r'
	case chess.BlackQueen:
		return 'q'
	case chess.BlackKing:
		return 'k'
	}
	return '.'
}

// Render writes an ASCII board representation of g to w.
// Ranks are displayed 8 (top) to 1 (bottom); files a-h left to right.
// Write errors are intentionally ignored: if the writer fails (e.g. closed pipe),
// partial output is acceptable for a terminal renderer.
func Render(g chess.Game, w io.Writer) {
	s := g.State

	_, _ = fmt.Fprintln(w, "  +---+---+---+---+---+---+---+---+")
	for rank := 7; rank >= 0; rank-- {
		_, _ = fmt.Fprintf(w, "%d |", rank+1)
		for file := 0; file <= 7; file++ {
			p := s.Board[chess.SquareOf(file, rank)]
			_, _ = fmt.Fprintf(w, " %c |", pieceASCII(p))
		}
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, "  +---+---+---+---+---+---+---+---+")
	}
	_, _ = fmt.Fprintln(w, "    a   b   c   d   e   f   g   h")

	// Side to move.
	if s.ActiveColor == chess.White {
		_, _ = fmt.Fprintln(w, "White to move")
	} else {
		_, _ = fmt.Fprintln(w, "Black to move")
	}
}
