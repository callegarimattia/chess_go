// Package chess implements a zero-dependency chess rules library.
// It provides FEN parsing, legal move generation, and immutable game state.
package chess

// Color represents which side has the move.
type Color uint8

const (
	White Color = 0
	Black Color = 1
)

// Piece encodes both color and type in a single byte.
// NoPiece = 0; White pieces = 1-6; Black pieces = 7-12.
type Piece uint8

const (
	NoPiece     Piece = 0
	WhitePawn   Piece = 1
	WhiteKnight Piece = 2
	WhiteBishop Piece = 3
	WhiteRook   Piece = 4
	WhiteQueen  Piece = 5
	WhiteKing   Piece = 6
	BlackPawn   Piece = 7
	BlackKnight Piece = 8
	BlackBishop Piece = 9
	BlackRook   Piece = 10
	BlackQueen  Piece = 11
	BlackKing   Piece = 12
)

// Square is a board index from 0 (a1) to 63 (h8).
// a1=0, b1=1, ..., h1=7, a2=8, ..., h8=63
type Square uint8

const NoSquare Square = 64

// Square constants for common squares.
const (
	A1 Square = 0
	B1 Square = 1
	C1 Square = 2
	D1 Square = 3
	E1 Square = 4
	F1 Square = 5
	G1 Square = 6
	H1 Square = 7

	A2 Square = 8
	B2 Square = 9
	C2 Square = 10
	D2 Square = 11
	E2 Square = 12
	F2 Square = 13
	G2 Square = 14
	H2 Square = 15

	A3 Square = 16
	B3 Square = 17
	C3 Square = 18
	D3 Square = 19
	E3 Square = 20
	F3 Square = 21
	G3 Square = 22
	H3 Square = 23

	A4 Square = 24
	B4 Square = 25
	C4 Square = 26
	D4 Square = 27
	E4 Square = 28
	F4 Square = 29
	G4 Square = 30
	H4 Square = 31

	A5 Square = 32
	B5 Square = 33
	C5 Square = 34
	D5 Square = 35
	E5 Square = 36
	F5 Square = 37
	G5 Square = 38
	H5 Square = 39

	A6 Square = 40
	B6 Square = 41
	C6 Square = 42
	D6 Square = 43
	E6 Square = 44
	F6 Square = 45
	G6 Square = 46
	H6 Square = 47

	A7 Square = 48
	B7 Square = 49
	C7 Square = 50
	D7 Square = 51
	E7 Square = 52
	F7 Square = 53
	G7 Square = 54
	H7 Square = 55

	A8 Square = 56
	B8 Square = 57
	C8 Square = 58
	D8 Square = 59
	E8 Square = 60
	F8 Square = 61
	G8 Square = 62
	H8 Square = 63
)

// Rank returns the rank of a square (0 = rank 1, 7 = rank 8).
func (s Square) Rank() int { return int(s) / 8 }

// File returns the file of a square (0 = file a, 7 = file h).
func (s Square) File() int { return int(s) % 8 }

// SquareOf constructs a Square from file (0-7) and rank (0-7).
func SquareOf(file, rank int) Square { return Square(rank*8 + file) }

// CastlingRight is a bitmask of available castling options.
type CastlingRight uint8

const (
	CastleWhiteKingside  CastlingRight = 0b0001
	CastleWhiteQueenside CastlingRight = 0b0010
	CastleBlackKingside  CastlingRight = 0b0100
	CastleBlackQueenside CastlingRight = 0b1000
	NoCastling           CastlingRight = 0b0000
	AllCastling          CastlingRight = 0b1111
)

// pieceColor returns the color of a piece (undefined for NoPiece).
func pieceColor(p Piece) Color {
	if p >= BlackPawn {
		return Black
	}
	return White
}

// pieceSymbol returns the ASCII letter for a piece (uppercase=White, lowercase=Black).
func pieceSymbol(p Piece) byte {
	switch p {
	case WhitePawn:
		return 'P'
	case WhiteKnight:
		return 'N'
	case WhiteBishop:
		return 'B'
	case WhiteRook:
		return 'R'
	case WhiteQueen:
		return 'Q'
	case WhiteKing:
		return 'K'
	case BlackPawn:
		return 'p'
	case BlackKnight:
		return 'n'
	case BlackBishop:
		return 'b'
	case BlackRook:
		return 'r'
	case BlackQueen:
		return 'q'
	case BlackKing:
		return 'k'
	}
	return '.'
}

// pieceFromSymbol converts an ASCII letter to a Piece.
func pieceFromSymbol(ch byte) Piece {
	switch ch {
	case 'P':
		return WhitePawn
	case 'N':
		return WhiteKnight
	case 'B':
		return WhiteBishop
	case 'R':
		return WhiteRook
	case 'Q':
		return WhiteQueen
	case 'K':
		return WhiteKing
	case 'p':
		return BlackPawn
	case 'n':
		return BlackKnight
	case 'b':
		return BlackBishop
	case 'r':
		return BlackRook
	case 'q':
		return BlackQueen
	case 'k':
		return BlackKing
	}
	return NoPiece
}
