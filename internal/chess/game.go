package chess

import "errors"

// ErrIllegalMove is returned by Game.Apply when the move is not in the legal move list.
var ErrIllegalMove = errors.New("illegal move")

// GameResult encodes terminal game states.
type GameResult uint8

const (
	InProgress               GameResult = 0
	WhiteWins                GameResult = 1
	BlackWins                GameResult = 2
	Stalemate                GameResult = 3
	DrawFiftyMove            GameResult = 4
	DrawThreefoldRepetition  GameResult = 5
	DrawInsufficientMaterial GameResult = 6
)

// GameState holds the complete chess position at a single point in time.
// It is a pure value type — copying is cheap (~80 bytes total).
type GameState struct {
	Board          [64]Piece
	ActiveColor    Color
	CastlingRights CastlingRight
	EnPassantSq    Square
	HalfMoveClock  uint8
	FullMoveNumber uint16
}

// Game wraps GameState with move history for draw detection and PGN export.
// Apply returns a new Game; the original is never mutated.
type Game struct {
	State   GameState
	history []GameState
}

// LegalMoves returns all legal moves from the current position.
func (g Game) LegalMoves() []Move {
	return generateLegalMoves(g.State)
}

// Apply applies move m to the game and returns a new Game.
// Returns ErrIllegalMove if m is not in LegalMoves().
func (g Game) Apply(m Move) (Game, error) {
	legal := g.LegalMoves()
	found := false
	for _, lm := range legal {
		if lm == m {
			found = true
			break
		}
	}
	if !found {
		return g, ErrIllegalMove
	}

	newState := applyMove(g.State, m)
	newHistory := make([]GameState, len(g.history)+1)
	copy(newHistory, g.history)
	newHistory[len(g.history)] = g.State

	return Game{
		State:   newState,
		history: newHistory,
	}, nil
}

// InCheck returns true if the active color's king is currently in check.
func (g Game) InCheck() bool {
	return isInCheck(g.State, g.State.ActiveColor)
}

// Result returns the current game result.
func (g Game) Result() GameResult {
	return detectResult(g.State, g.history)
}

// ToFEN returns the FEN string for the current game state.
func (g Game) ToFEN() string {
	return stateToFEN(g.State)
}

// ToPGN returns a minimal PGN string for the game.
func (g Game) ToPGN() string {
	return buildPGN(g)
}

// applyMove applies a move to a GameState and returns the new state.
// The move is assumed to be legal.
func applyMove(s GameState, m Move) GameState {
	ns := s // copy

	movingPiece := ns.Board[m.From]
	capturedPiece := ns.Board[m.To]

	// En passant capture: remove the captured pawn.
	isEnPassant := false
	if (movingPiece == WhitePawn || movingPiece == BlackPawn) &&
		m.To == ns.EnPassantSq && ns.EnPassantSq != NoSquare {
		isEnPassant = true
		if movingPiece == WhitePawn {
			// Captured pawn is one rank below the target square.
			ns.Board[Square(m.To-8)] = NoPiece
		} else {
			ns.Board[Square(m.To+8)] = NoPiece
		}
	}

	// Move the piece.
	ns.Board[m.From] = NoPiece
	if m.Promotion != NoPiece {
		ns.Board[m.To] = m.Promotion
	} else {
		ns.Board[m.To] = movingPiece
	}

	// Castling: move the rook.
	if movingPiece == WhiteKing && m.From == E1 {
		switch m.To {
		case G1:
			ns.Board[H1] = NoPiece
			ns.Board[F1] = WhiteRook
		case C1:
			ns.Board[A1] = NoPiece
			ns.Board[D1] = WhiteRook
		}
	}
	if movingPiece == BlackKing && m.From == E8 {
		switch m.To {
		case G8:
			ns.Board[H8] = NoPiece
			ns.Board[F8] = BlackRook
		case C8:
			ns.Board[A8] = NoPiece
			ns.Board[D8] = BlackRook
		}
	}

	// Update castling rights.
	if movingPiece == WhiteKing {
		ns.CastlingRights &^= CastleWhiteKingside | CastleWhiteQueenside
	}
	if movingPiece == BlackKing {
		ns.CastlingRights &^= CastleBlackKingside | CastleBlackQueenside
	}
	if m.From == A1 || m.To == A1 {
		ns.CastlingRights &^= CastleWhiteQueenside
	}
	if m.From == H1 || m.To == H1 {
		ns.CastlingRights &^= CastleWhiteKingside
	}
	if m.From == A8 || m.To == A8 {
		ns.CastlingRights &^= CastleBlackQueenside
	}
	if m.From == H8 || m.To == H8 {
		ns.CastlingRights &^= CastleBlackKingside
	}

	// Update en passant square.
	ns.EnPassantSq = NoSquare
	if movingPiece == WhitePawn && m.To.Rank()-m.From.Rank() == 2 {
		ns.EnPassantSq = Square(m.From + 8)
	}
	if movingPiece == BlackPawn && m.From.Rank()-m.To.Rank() == 2 {
		ns.EnPassantSq = Square(m.From - 8)
	}

	// Update half-move clock.
	if movingPiece == WhitePawn || movingPiece == BlackPawn ||
		capturedPiece != NoPiece || isEnPassant {
		ns.HalfMoveClock = 0
	} else {
		ns.HalfMoveClock++
	}

	// Update full-move number.
	if s.ActiveColor == Black {
		ns.FullMoveNumber++
	}

	// Switch active color.
	if s.ActiveColor == White {
		ns.ActiveColor = Black
	} else {
		ns.ActiveColor = White
	}

	return ns
}

// isSquareAttackedBy returns true if the given square is attacked by any piece of the given color.
func isSquareAttackedBy(s GameState, sq Square, byColor Color) bool {
	// Check pawns.
	if byColor == White {
		// White pawns attack diagonally upward.
		if sq.Rank() > 0 {
			if sq.File() > 0 {
				attacker := Square(sq - 9)
				if s.Board[attacker] == WhitePawn {
					return true
				}
			}
			if sq.File() < 7 {
				attacker := Square(sq - 7)
				if s.Board[attacker] == WhitePawn {
					return true
				}
			}
		}
	} else {
		// Black pawns attack diagonally downward.
		if sq.Rank() < 7 {
			if sq.File() > 0 {
				attacker := Square(sq + 7)
				if s.Board[attacker] == BlackPawn {
					return true
				}
			}
			if sq.File() < 7 {
				attacker := Square(sq + 9)
				if s.Board[attacker] == BlackPawn {
					return true
				}
			}
		}
	}

	// Knights.
	knightPiece := WhiteKnight
	if byColor == Black {
		knightPiece = BlackKnight
	}
	knightOffsets := [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	for _, off := range knightOffsets {
		r := sq.Rank() + off[0]
		f := sq.File() + off[1]
		if r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			if s.Board[SquareOf(f, r)] == knightPiece {
				return true
			}
		}
	}

	// Sliding pieces (rook, queen for rank/file; bishop, queen for diagonals).
	rookPiece := WhiteRook
	queenPiece := WhiteQueen
	bishopPiece := WhiteBishop
	kingPiece := WhiteKing
	if byColor == Black {
		rookPiece = BlackRook
		queenPiece = BlackQueen
		bishopPiece = BlackBishop
		kingPiece = BlackKing
	}

	// Rook/Queen directions (rank and file).
	rookDirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for _, dir := range rookDirs {
		r := sq.Rank() + dir[0]
		f := sq.File() + dir[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := s.Board[SquareOf(f, r)]
			if p != NoPiece {
				if p == rookPiece || p == queenPiece {
					return true
				}
				break
			}
			r += dir[0]
			f += dir[1]
		}
	}

	// Bishop/Queen directions (diagonals).
	bishopDirs := [4][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, dir := range bishopDirs {
		r := sq.Rank() + dir[0]
		f := sq.File() + dir[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := s.Board[SquareOf(f, r)]
			if p != NoPiece {
				if p == bishopPiece || p == queenPiece {
					return true
				}
				break
			}
			r += dir[0]
			f += dir[1]
		}
	}

	// King (one square in any direction).
	for dr := -1; dr <= 1; dr++ {
		for df := -1; df <= 1; df++ {
			if dr == 0 && df == 0 {
				continue
			}
			r := sq.Rank() + dr
			f := sq.File() + df
			if r >= 0 && r <= 7 && f >= 0 && f <= 7 {
				if s.Board[SquareOf(f, r)] == kingPiece {
					return true
				}
			}
		}
	}

	return false
}

// isInCheck returns true if the given color's king is in check.
func isInCheck(s GameState, color Color) bool {
	// Find the king.
	kingPiece := WhiteKing
	if color == Black {
		kingPiece = BlackKing
	}
	var kingSq Square
	for sq := Square(0); sq < 64; sq++ {
		if s.Board[sq] == kingPiece {
			kingSq = sq
			break
		}
	}

	// Check if the king square is attacked by the opponent.
	var opponent Color
	if color == White {
		opponent = Black
	} else {
		opponent = White
	}
	return isSquareAttackedBy(s, kingSq, opponent)
}

// detectResult returns the current game result.
func detectResult(s GameState, history []GameState) GameResult {
	// Fifty-move rule.
	if s.HalfMoveClock >= 100 {
		return DrawFiftyMove
	}

	// Insufficient material (kings only).
	if isInsufficientMaterial(s) {
		return DrawInsufficientMaterial
	}

	// Threefold repetition.
	if isThreefoldRepetition(s, history) {
		return DrawThreefoldRepetition
	}

	// Check for legal moves.
	moves := generateLegalMoves(s)
	if len(moves) > 0 {
		return InProgress
	}

	// No legal moves: checkmate or stalemate.
	if isInCheck(s, s.ActiveColor) {
		if s.ActiveColor == White {
			return BlackWins
		}
		return WhiteWins
	}
	return Stalemate
}

// isInsufficientMaterial returns true if neither side can deliver checkmate.
func isInsufficientMaterial(s GameState) bool {
	for sq := Square(0); sq < 64; sq++ {
		p := s.Board[sq]
		if p == NoPiece || p == WhiteKing || p == BlackKing {
			continue
		}
		// Any pawn, rook, or queen = sufficient material.
		if p == WhitePawn || p == BlackPawn ||
			p == WhiteRook || p == BlackRook ||
			p == WhiteQueen || p == BlackQueen {
			return false
		}
		// Two knights vs lone king is technically a draw but very rare.
		// For the skeleton we only handle the trivial king-vs-king case.
		return false
	}
	return true
}

// isThreefoldRepetition returns true if the current position has appeared at least 3 times.
func isThreefoldRepetition(current GameState, history []GameState) bool {
	count := 1
	for _, h := range history {
		if positionsEqual(current, h) {
			count++
			if count >= 3 {
				return true
			}
		}
	}
	return false
}

// positionsEqual compares the position-relevant fields of two GameStates.
func positionsEqual(a, b GameState) bool {
	return a.Board == b.Board &&
		a.ActiveColor == b.ActiveColor &&
		a.CastlingRights == b.CastlingRights &&
		a.EnPassantSq == b.EnPassantSq
}
