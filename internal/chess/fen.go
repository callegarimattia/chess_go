package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidFEN is returned by NewGameFromFEN when the FEN string is malformed.
var ErrInvalidFEN = errors.New("invalid FEN")

// NewGameFromFEN parses a FEN string and returns a new Game.
// Returns ErrInvalidFEN if the string is malformed or invalid.
func NewGameFromFEN(fen string) (Game, error) {
	if fen == "" {
		return Game{}, ErrInvalidFEN
	}

	parts := strings.Fields(fen)
	if len(parts) != 6 {
		return Game{}, ErrInvalidFEN
	}

	var state GameState

	// Parse piece placement.
	board, err := parsePiecePlacement(parts[0])
	if err != nil {
		return Game{}, err
	}
	state.Board = board

	// Parse active color.
	switch parts[1] {
	case "w":
		state.ActiveColor = White
	case "b":
		state.ActiveColor = Black
	default:
		return Game{}, ErrInvalidFEN
	}

	// Parse castling rights.
	cr, err := parseCastlingRights(parts[2])
	if err != nil {
		return Game{}, err
	}
	state.CastlingRights = cr

	// Parse en passant square.
	ep, err := parseEnPassantSquare(parts[3])
	if err != nil {
		return Game{}, err
	}
	state.EnPassantSq = ep

	// Parse half-move clock.
	hmc, err := strconv.ParseUint(parts[4], 10, 8)
	if err != nil {
		return Game{}, ErrInvalidFEN
	}
	state.HalfMoveClock = uint8(hmc)

	// Parse full-move number.
	fmn, err := strconv.ParseUint(parts[5], 10, 16)
	if err != nil || fmn == 0 {
		return Game{}, ErrInvalidFEN
	}
	state.FullMoveNumber = uint16(fmn)

	return Game{State: state}, nil
}

// parsePiecePlacement parses the piece placement field of a FEN string.
func parsePiecePlacement(s string) ([64]Piece, error) {
	var board [64]Piece
	ranks := strings.Split(s, "/")
	if len(ranks) != 8 {
		return board, ErrInvalidFEN
	}

	// FEN ranks are ordered 8 to 1 (rank index 7 to 0).
	for rankIdx, rankStr := range ranks {
		rank := 7 - rankIdx
		file := 0
		for i := 0; i < len(rankStr); i++ {
			ch := rankStr[i]
			if ch >= '1' && ch <= '8' {
				file += int(ch - '0')
			} else {
				p := pieceFromSymbol(ch)
				if p == NoPiece {
					return board, fmt.Errorf("%w: invalid piece character %q", ErrInvalidFEN, ch)
				}
				if file > 7 {
					return board, ErrInvalidFEN
				}
				board[SquareOf(file, rank)] = p
				file++
			}
		}
		if file != 8 {
			return board, ErrInvalidFEN
		}
	}
	return board, nil
}

// parseCastlingRights parses the castling availability field.
func parseCastlingRights(s string) (CastlingRight, error) {
	if s == "-" {
		return NoCastling, nil
	}
	var cr CastlingRight
	for _, ch := range s {
		switch ch {
		case 'K':
			cr |= CastleWhiteKingside
		case 'Q':
			cr |= CastleWhiteQueenside
		case 'k':
			cr |= CastleBlackKingside
		case 'q':
			cr |= CastleBlackQueenside
		default:
			return NoCastling, ErrInvalidFEN
		}
	}
	return cr, nil
}

// parseEnPassantSquare parses the en passant target square field.
func parseEnPassantSquare(s string) (Square, error) {
	if s == "-" {
		return NoSquare, nil
	}
	if len(s) != 2 {
		return NoSquare, ErrInvalidFEN
	}
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return NoSquare, ErrInvalidFEN
	}
	return SquareOf(file, rank), nil
}

// stateToFEN converts a GameState to its FEN string representation.
func stateToFEN(s GameState) string {
	var sb strings.Builder

	// Piece placement.
	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := 0; file <= 7; file++ {
			p := s.Board[SquareOf(file, rank)]
			if p == NoPiece {
				empty++
			} else {
				if empty > 0 {
					sb.WriteByte(byte('0' + empty))
					empty = 0
				}
				sb.WriteByte(pieceSymbol(p))
			}
		}
		if empty > 0 {
			sb.WriteByte(byte('0' + empty))
		}
		if rank > 0 {
			sb.WriteByte('/')
		}
	}

	sb.WriteByte(' ')

	// Active color.
	if s.ActiveColor == White {
		sb.WriteByte('w')
	} else {
		sb.WriteByte('b')
	}

	sb.WriteByte(' ')

	// Castling rights.
	if s.CastlingRights == NoCastling {
		sb.WriteByte('-')
	} else {
		if s.CastlingRights&CastleWhiteKingside != 0 {
			sb.WriteByte('K')
		}
		if s.CastlingRights&CastleWhiteQueenside != 0 {
			sb.WriteByte('Q')
		}
		if s.CastlingRights&CastleBlackKingside != 0 {
			sb.WriteByte('k')
		}
		if s.CastlingRights&CastleBlackQueenside != 0 {
			sb.WriteByte('q')
		}
	}

	sb.WriteByte(' ')

	// En passant square.
	if s.EnPassantSq == NoSquare {
		sb.WriteByte('-')
	} else {
		sb.WriteByte(byte('a' + s.EnPassantSq.File()))
		sb.WriteByte(byte('1' + s.EnPassantSq.Rank()))
	}

	sb.WriteByte(' ')
	fmt.Fprintf(&sb, "%d %d", s.HalfMoveClock, s.FullMoveNumber)

	return sb.String()
}
