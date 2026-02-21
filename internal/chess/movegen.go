package chess

// generateLegalMoves generates all legal moves for the active color.
// It generates pseudo-legal moves then filters out those that leave the king in check.
func generateLegalMoves(s GameState) []Move {
	pseudo := generatePseudoLegal(s)
	legal := pseudo[:0:len(pseudo)]
	legal = legal[:0]

	for _, m := range pseudo {
		ns := applyMove(s, m)
		// After the move, the side that just moved must not be in check.
		if !isInCheck(ns, s.ActiveColor) {
			legal = append(legal, m)
		}
	}
	return legal
}

// generatePseudoLegal generates all pseudo-legal moves (may leave king in check).
func generatePseudoLegal(s GameState) []Move {
	var moves []Move

	for sq := Square(0); sq < 64; sq++ {
		p := s.Board[sq]
		if p == NoPiece {
			continue
		}
		if pieceColor(p) != s.ActiveColor {
			continue
		}

		switch p {
		case WhitePawn:
			moves = append(moves, genWhitePawnMoves(s, sq)...)
		case BlackPawn:
			moves = append(moves, genBlackPawnMoves(s, sq)...)
		case WhiteKnight, BlackKnight:
			moves = append(moves, genKnightMoves(s, sq, p)...)
		case WhiteBishop, BlackBishop:
			moves = append(moves, genSlidingMoves(s, sq, p, [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}})...)
		case WhiteRook, BlackRook:
			moves = append(moves, genSlidingMoves(s, sq, p, [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}})...)
		case WhiteQueen, BlackQueen:
			moves = append(moves, genSlidingMoves(s, sq, p, [][2]int{
				{0, 1}, {0, -1}, {1, 0}, {-1, 0},
				{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
			})...)
		case WhiteKing, BlackKing:
			moves = append(moves, genKingMoves(s, sq, p)...)
		}
	}

	return moves
}

// genWhitePawnMoves generates pseudo-legal moves for a White pawn.
func genWhitePawnMoves(s GameState, sq Square) []Move {
	var moves []Move
	rank := sq.Rank()
	file := sq.File()

	// Single push.
	if rank < 7 {
		to := Square(sq + 8)
		if s.Board[to] == NoPiece {
			if rank == 6 {
				// Promotion.
				moves = append(moves, promotionMoves(sq, to, White)...)
			} else {
				moves = append(moves, Move{From: sq, To: to})
				// Double push from rank 2 (index 1).
				if rank == 1 {
					to2 := Square(sq + 16)
					if s.Board[to2] == NoPiece {
						moves = append(moves, Move{From: sq, To: to2})
					}
				}
			}
		}
	}

	// Captures.
	if rank < 7 {
		if file > 0 {
			to := Square(sq + 7)
			if s.Board[to] != NoPiece && pieceColor(s.Board[to]) == Black {
				if rank == 6 {
					moves = append(moves, promotionMoves(sq, to, White)...)
				} else {
					moves = append(moves, Move{From: sq, To: to})
				}
			}
			// En passant.
			if s.EnPassantSq != NoSquare && to == s.EnPassantSq {
				moves = append(moves, Move{From: sq, To: to})
			}
		}
		if file < 7 {
			to := Square(sq + 9)
			if s.Board[to] != NoPiece && pieceColor(s.Board[to]) == Black {
				if rank == 6 {
					moves = append(moves, promotionMoves(sq, to, White)...)
				} else {
					moves = append(moves, Move{From: sq, To: to})
				}
			}
			// En passant.
			if s.EnPassantSq != NoSquare && to == s.EnPassantSq {
				moves = append(moves, Move{From: sq, To: to})
			}
		}
	}

	return moves
}

// genBlackPawnMoves generates pseudo-legal moves for a Black pawn.
func genBlackPawnMoves(s GameState, sq Square) []Move {
	var moves []Move
	rank := sq.Rank()
	file := sq.File()

	// Single push.
	if rank > 0 {
		to := Square(sq - 8)
		if s.Board[to] == NoPiece {
			if rank == 1 {
				// Promotion.
				moves = append(moves, promotionMoves(sq, to, Black)...)
			} else {
				moves = append(moves, Move{From: sq, To: to})
				// Double push from rank 7 (index 6).
				if rank == 6 {
					to2 := Square(sq - 16)
					if s.Board[to2] == NoPiece {
						moves = append(moves, Move{From: sq, To: to2})
					}
				}
			}
		}
	}

	// Captures.
	if rank > 0 {
		if file > 0 {
			to := Square(sq - 9)
			if s.Board[to] != NoPiece && pieceColor(s.Board[to]) == White {
				if rank == 1 {
					moves = append(moves, promotionMoves(sq, to, Black)...)
				} else {
					moves = append(moves, Move{From: sq, To: to})
				}
			}
			// En passant.
			if s.EnPassantSq != NoSquare && to == s.EnPassantSq {
				moves = append(moves, Move{From: sq, To: to})
			}
		}
		if file < 7 {
			to := Square(sq - 7)
			if s.Board[to] != NoPiece && pieceColor(s.Board[to]) == White {
				if rank == 1 {
					moves = append(moves, promotionMoves(sq, to, Black)...)
				} else {
					moves = append(moves, Move{From: sq, To: to})
				}
			}
			// En passant.
			if s.EnPassantSq != NoSquare && to == s.EnPassantSq {
				moves = append(moves, Move{From: sq, To: to})
			}
		}
	}

	return moves
}

// promotionMoves returns the four promotion moves for a pawn reaching the last rank.
func promotionMoves(from, to Square, color Color) []Move {
	if color == White {
		return []Move{
			{From: from, To: to, Promotion: WhiteQueen},
			{From: from, To: to, Promotion: WhiteRook},
			{From: from, To: to, Promotion: WhiteBishop},
			{From: from, To: to, Promotion: WhiteKnight},
		}
	}
	return []Move{
		{From: from, To: to, Promotion: BlackQueen},
		{From: from, To: to, Promotion: BlackRook},
		{From: from, To: to, Promotion: BlackBishop},
		{From: from, To: to, Promotion: BlackKnight},
	}
}

// genKnightMoves generates pseudo-legal knight moves.
func genKnightMoves(s GameState, sq Square, p Piece) []Move {
	var moves []Move
	color := pieceColor(p)
	offsets := [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	for _, off := range offsets {
		r := sq.Rank() + off[0]
		f := sq.File() + off[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}
		to := SquareOf(f, r)
		dest := s.Board[to]
		if dest == NoPiece || pieceColor(dest) != color {
			moves = append(moves, Move{From: sq, To: to})
		}
	}
	return moves
}

// genSlidingMoves generates pseudo-legal moves for sliding pieces (rook, bishop, queen).
func genSlidingMoves(s GameState, sq Square, p Piece, dirs [][2]int) []Move {
	var moves []Move
	color := pieceColor(p)
	for _, dir := range dirs {
		r := sq.Rank() + dir[0]
		f := sq.File() + dir[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			to := SquareOf(f, r)
			dest := s.Board[to]
			if dest == NoPiece {
				moves = append(moves, Move{From: sq, To: to})
			} else if pieceColor(dest) != color {
				moves = append(moves, Move{From: sq, To: to})
				break
			} else {
				break
			}
			r += dir[0]
			f += dir[1]
		}
	}
	return moves
}

// genKingMoves generates pseudo-legal king moves including castling.
func genKingMoves(s GameState, sq Square, p Piece) []Move {
	var moves []Move
	color := pieceColor(p)

	// Normal one-square moves.
	for dr := -1; dr <= 1; dr++ {
		for df := -1; df <= 1; df++ {
			if dr == 0 && df == 0 {
				continue
			}
			r := sq.Rank() + dr
			f := sq.File() + df
			if r < 0 || r > 7 || f < 0 || f > 7 {
				continue
			}
			to := SquareOf(f, r)
			dest := s.Board[to]
			if dest == NoPiece || pieceColor(dest) != color {
				moves = append(moves, Move{From: sq, To: to})
			}
		}
	}

	// Castling.
	opponent := Black
	if color == Black {
		opponent = White
	}

	if color == White && sq == E1 && !isSquareAttackedBy(s, E1, opponent) {
		// Kingside: e1-f1-g1 must be empty, f1 not attacked.
		if s.CastlingRights&CastleWhiteKingside != 0 &&
			s.Board[F1] == NoPiece && s.Board[G1] == NoPiece &&
			!isSquareAttackedBy(s, F1, opponent) {
			moves = append(moves, Move{From: E1, To: G1})
		}
		// Queenside: e1-d1-c1 must be empty, d1 not attacked.
		if s.CastlingRights&CastleWhiteQueenside != 0 &&
			s.Board[D1] == NoPiece && s.Board[C1] == NoPiece && s.Board[B1] == NoPiece &&
			!isSquareAttackedBy(s, D1, opponent) {
			moves = append(moves, Move{From: E1, To: C1})
		}
	}

	if color == Black && sq == E8 && !isSquareAttackedBy(s, E8, opponent) {
		// Kingside.
		if s.CastlingRights&CastleBlackKingside != 0 &&
			s.Board[F8] == NoPiece && s.Board[G8] == NoPiece &&
			!isSquareAttackedBy(s, F8, opponent) {
			moves = append(moves, Move{From: E8, To: G8})
		}
		// Queenside.
		if s.CastlingRights&CastleBlackQueenside != 0 &&
			s.Board[D8] == NoPiece && s.Board[C8] == NoPiece && s.Board[B8] == NoPiece &&
			!isSquareAttackedBy(s, D8, opponent) {
			moves = append(moves, Move{From: E8, To: C8})
		}
	}

	return moves
}
