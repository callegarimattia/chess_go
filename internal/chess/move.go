package chess

import "fmt"

// Move represents a chess move as a (from, to, promotion) triple.
// It is a value type — no heap allocation in the hot search path.
type Move struct {
	From      Square
	To        Square
	Promotion Piece // NoPiece if not a promotion
}

// UCIString returns the UCI representation of the move (e.g. "e2e4", "e7e8q").
func (m Move) UCIString() string {
	fromFile := rune('a' + m.From.File())
	fromRank := rune('1' + m.From.Rank())
	toFile := rune('a' + m.To.File())
	toRank := rune('1' + m.To.Rank())

	s := fmt.Sprintf("%c%c%c%c", fromFile, fromRank, toFile, toRank)
	if m.Promotion != NoPiece {
		promoSymbol := pieceSymbol(m.Promotion)
		// UCI uses lowercase for promotion piece
		if promoSymbol >= 'A' && promoSymbol <= 'Z' {
			promoSymbol += 32
		}
		s += string(promoSymbol)
	}
	return s
}

// IsPromotion returns true if this move includes a promotion piece.
func (m Move) IsPromotion() bool { return m.Promotion != NoPiece }

// IsCastle returns true if this move is a king castling move (king moves 2 squares).
func (m Move) IsCastle() bool {
	fileDiff := m.To.File() - m.From.File()
	if fileDiff < 0 {
		fileDiff = -fileDiff
	}
	return fileDiff == 2
}
