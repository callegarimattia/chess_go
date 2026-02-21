// chess_steps_test.go — Executable specifications for the chess logic package.
//
// Mirrors: milestone-1-chess-package.feature
// Driving port: internal/chess public API
//   - chess.NewGameFromFEN(fen string) (chess.Game, error)
//   - chess.Game.LegalMoves() []chess.Move
//   - chess.Game.Apply(m chess.Move) (chess.Game, error)
//   - chess.Game.InCheck() bool
//   - chess.Game.Result() chess.GameResult
//   - chess.Game.ToFEN() string
//   - chess.Game.ToPGN() string
//   - chess.Move.UCIString() string
//   - chess.Move.SANString(g chess.Game) string
//
// CM-A compliance: all imports point to internal/chess only.
// CM-B compliance: test names use domain language exclusively.

package acceptance_test

import (
	"strings"
	"testing"

	chess "chess_go/internal/chess"
)

// ─── Walking Skeleton ─────────────────────────────────────────────────────────

// TestWalkingSkeleton_PipelineConnectsAllLayers validates US-00 / AC-00.
// Gherkin: walking-skeleton.feature — "Walking skeleton pipeline connects all architectural layers"
func TestWalkingSkeleton_PipelineConnectsAllLayers(t *testing.T) {
	_ = requiresProduction("internal/chess", "internal/engine", "internal/tui")

	// Phase 1: FEN parse.
	game, err := chess.NewGameFromFEN(StartingFEN)
	if err != nil {
		t.Fatalf("FEN parse failed: %v", err)
	}

	// Phase 2: Legal move generation.
	moves := game.LegalMoves()
	if len(moves) != 20 {
		t.Fatalf("expected 20 legal moves from start, got %d", len(moves))
	}

	// Phase 3: Random selection — exercised via engine package (see engine_steps_test.go).

	// Phase 4: Apply move.
	newGame, err := game.Apply(moves[0])
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Phase 5: Render — exercised via tui package (see tui_steps_test.go).

	// Original state unchanged (immutability).
	if game.ToFEN() != StartingFEN {
		t.Errorf("original game mutated: got %s", game.ToFEN())
	}
	_ = newGame

	t.Log("walking skeleton: pipeline from FEN parse to Apply is connected")
}

// ─── FEN Parsing ─────────────────────────────────────────────────────────────

// TestFENParser_ValidStartingPosition validates US-01 / AC-01-01.
// Gherkin: "Library consumer parses the starting position from a FEN string"
func TestFENParser_ValidStartingPosition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// require.NotZero(t, game)
	// assert.Equal(t, chess.White, game.State.ActiveColor)
	// assert.Equal(t, chess.AllCastling, game.State.CastlingRights)
	// assert.Equal(t, chess.NoSquare, game.State.EnPassantSq)
	// assert.Equal(t, uint8(0), game.State.HalfMoveClock)
	// assert.Equal(t, uint16(1), game.State.FullMoveNumber)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestFENParser_MalformedFENReturnsTypedError validates US-01 / AC-01-02.
// Gherkin: "Library consumer receives a typed error for a malformed FEN string"
func TestFENParser_MalformedFENReturnsTypedError(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	cases := []struct {
		name string
		fen  string
	}{
		{"plain garbage", "not-a-valid-fen"},
		{"too few fields", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq"},
		{"invalid piece char", "rnbqkxnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"empty string", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// _, err := chess.NewGameFromFEN(tc.fen)
			// require.Error(t, err)
			// require.True(t, errors.Is(err, chess.ErrInvalidFEN),
			//     "expected ErrInvalidFEN, got %T: %v", err, err)
		})
	}
	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestFENParser_RoundTrip validates US-01 / AC-01-03.
// Gherkin: "Game state serialises back to the same FEN string it was loaded from"
func TestFENParser_RoundTrip(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	fens := []string{
		StartingFEN,
		AfterE2E4FEN,
		FoolsMateFEN,
		KiwipeteFEN,
	}
	for _, fen := range fens {
		t.Run(fen, func(t *testing.T) {
			// game, err := chess.NewGameFromFEN(fen)
			// require.NoError(t, err)
			// got := game.ToFEN()
			// assert.Equal(t, fen, got, "FEN round-trip failed")
		})
	}
	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestFENParser_EnPassantPreservedAfterPawnAdvance validates FEN round-trip after e2e4.
// Gherkin: "FEN round-trip preserves en passant square after a pawn double advance"
func TestFENParser_EnPassantPreservedAfterPawnAdvance(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e2e4")
	// game, err = game.Apply(move)
	// require.NoError(t, err)
	// fen := game.ToFEN()
	// assert.Contains(t, fen, "e3", "en passant square missing from FEN after e2e4")
	// reloaded, err := chess.NewGameFromFEN(fen)
	// require.NoError(t, err)
	// assert.Equal(t, game.State.EnPassantSq, reloaded.State.EnPassantSq)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Legal Move Generation ────────────────────────────────────────────────────

// TestLegalMoves_StartingPositionHas20Moves validates US-02 / AC-02-01.
// Gherkin: "Library consumer receives exactly 20 legal moves from the starting position"
func TestLegalMoves_StartingPositionHas20Moves(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// moves := game.LegalMoves()
	// assert.Len(t, moves, 20, "starting position must have exactly 20 legal moves")
	// for _, m := range moves {
	//     piece := game.State.Board[m.From]
	//     isPawn := piece == chess.WhitePawn
	//     isKnight := piece == chess.WhiteKnight
	//     assert.True(t, isPawn || isKnight,
	//         "expected only pawn/knight moves from start, got move %s from piece %v", m.UCIString(), piece)
	// }

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestLegalMoves_NoIllegalMovesReturned validates US-02 / AC-02-02.
// Gherkin: "Library consumer receives no illegal moves in any generated move list"
func TestLegalMoves_NoIllegalMovesReturned(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	positions := []string{
		StartingFEN,
		KiwipeteFEN,
		AfterE2E4FEN,
	}
	for _, fen := range positions {
		t.Run(fen[:30], func(t *testing.T) {
			// game, err := chess.NewGameFromFEN(fen)
			// require.NoError(t, err)
			// for _, m := range game.LegalMoves() {
			//     applied, err := game.Apply(m)
			//     require.NoError(t, err, "Apply of LegalMoves move must not fail: %s", m.UCIString())
			//     // After applying the move, switch perspective and check the previous mover is not in check
			//     // (InCheck on the *new* state checks the *new* active colour)
			//     _ = applied
			// }
		})
	}
	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestLegalMoves_CheckmatePositionReturnsEmptyList validates US-02 / AC-02-03.
// Gherkin: "Library consumer receives an empty move list from a checkmate position"
func TestLegalMoves_CheckmatePositionReturnsEmptyList(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(FoolsMateFEN)
	// require.NoError(t, err)
	// moves := game.LegalMoves()
	// assert.Empty(t, moves, "checkmate position must return zero legal moves")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestLegalMoves_StalematePositionReturnsEmptyList validates US-02 / AC-02-04.
// Gherkin: "Library consumer receives an empty move list from a stalemate position"
func TestLegalMoves_StalematePositionReturnsEmptyList(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StalemateFEN)
	// require.NoError(t, err)
	// moves := game.LegalMoves()
	// assert.Empty(t, moves, "stalemate position must return zero legal moves")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Apply Move / Immutability ────────────────────────────────────────────────

// TestApplyMove_PawnDoubleAdvanceProducesCorrectState validates US-03 / AC-03-01.
// Gherkin: "Library consumer applies a pawn double advance and receives the correct new state"
func TestApplyMove_PawnDoubleAdvanceProducesCorrectState(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e2e4")
	//
	// newGame, err := game.Apply(move)
	// require.NoError(t, err)
	//
	// // New state assertions
	// assert.Equal(t, chess.Black, newGame.State.ActiveColor, "active color must switch to Black")
	// assert.Equal(t, chess.E3, newGame.State.EnPassantSq, "en passant square must be e3")
	// assert.Equal(t, chess.WhitePawn, newGame.State.Board[chess.E4], "pawn must be on e4")
	// assert.Equal(t, chess.NoPiece, newGame.State.Board[chess.E2], "e2 must be empty")
	//
	// // Original state unchanged
	// assert.Equal(t, chess.White, game.State.ActiveColor, "original game must still have White to move")
	// assert.Equal(t, chess.WhitePawn, game.State.Board[chess.E2], "original pawn must still be on e2")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestApplyMove_IllegalMoveReturnsTypedError validates US-03 / AC-03-02.
// Gherkin: "Library consumer receives a typed error when applying a move not in the legal move list"
func TestApplyMove_IllegalMoveReturnsTypedError(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// illegalMove := chess.Move{From: chess.E2, To: chess.E5}
	// _, err = game.Apply(illegalMove)
	// require.Error(t, err)
	// assert.True(t, errors.Is(err, chess.ErrIllegalMove))
	// assert.Equal(t, StartingFEN, game.ToFEN(), "original state must be unchanged after illegal Apply")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestApplyMove_SequenceProducesIndependentStates validates immutability across a move sequence.
// Gherkin: "Library consumer applies a sequence of moves and each new state is independent"
func TestApplyMove_SequenceProducesIndependentStates(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// gameA, err := game.Apply(mustParseUCI(t, game, "e2e4"))
	// require.NoError(t, err)
	// gameB, err := gameA.Apply(mustParseUCI(t, gameA, "e7e5"))
	// require.NoError(t, err)
	// gameC, err := gameB.Apply(mustParseUCI(t, gameB, "g1f3"))
	// require.NoError(t, err)
	//
	// assert.Equal(t, chess.Black, gameA.State.ActiveColor)
	// assert.Equal(t, chess.White, gameB.State.ActiveColor)
	// assert.Equal(t, chess.WhiteKnight, gameC.State.Board[chess.F3])
	// assert.Equal(t, chess.Black, gameC.State.ActiveColor)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Check and Checkmate Detection ───────────────────────────────────────────

// TestResult_CheckmateDetectedInFoolsMate validates US-04 / AC-04-01.
// Gherkin: "Library consumer detects checkmate in the Fool's Mate position"
func TestResult_CheckmateDetectedInFoolsMate(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(FoolsMateFEN)
	// require.NoError(t, err)
	// result := game.Result()
	// assert.Equal(t, chess.BlackWins, result, "expected BlackWins (checkmate), got %v", result)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestInCheck_TrueInFoolsMatePosition validates US-04 / AC-04-02.
// Gherkin: "Library consumer detects that the active king is in check"
func TestInCheck_TrueInFoolsMatePosition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(FoolsMateFEN)
	// require.NoError(t, err)
	// assert.True(t, game.InCheck(), "king must be in check in Fool's Mate position")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestInCheck_FalseInStartingPosition validates US-04 / AC-04-03.
// Gherkin: "Library consumer detects that the active king is not in check"
func TestInCheck_FalseInStartingPosition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// assert.False(t, game.InCheck())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestResult_InProgressForStartingPosition validates that a normal position is InProgress.
// Gherkin: "Library consumer detects that a game in progress has no terminal result"
func TestResult_InProgressForStartingPosition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// assert.Equal(t, chess.InProgress, game.Result())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Draw Detection ───────────────────────────────────────────────────────────

// TestResult_DrawByFiftyMoveRule validates US-05 / AC-05-01.
// Gherkin: "Library consumer detects a draw by the fifty-move rule"
func TestResult_DrawByFiftyMoveRule(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// Construct a position with half-move clock = 100 via FEN.
	// fen := "8/8/8/8/8/8/8/K6k w - - 100 101"
	// game, err := chess.NewGameFromFEN(fen)
	// require.NoError(t, err)
	// assert.Equal(t, chess.DrawFiftyMove, game.Result())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestResult_DrawByThreefoldRepetition validates US-05 / AC-05-02.
// Gherkin: "Library consumer detects a draw by threefold repetition"
func TestResult_DrawByThreefoldRepetition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// Shuffle knight g1-f3-g1-f3-g1 (starting position appears three times).
	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// moves := []string{"g1f3", "g8f6", "f3g1", "f6g8", "g1f3", "g8f6", "f3g1", "f6g8"}
	// for _, uci := range moves {
	//     m := mustParseUCI(t, game, uci)
	//     game, err = game.Apply(m)
	//     require.NoError(t, err)
	// }
	// assert.Equal(t, chess.DrawThreefoldRepetition, game.Result())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestResult_DrawByInsufficientMaterialKingsOnly validates US-05 / AC-05-03.
// Gherkin: "Library consumer detects a draw by insufficient material with kings only"
func TestResult_DrawByInsufficientMaterialKingsOnly(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(KingsOnlyFEN)
	// require.NoError(t, err)
	// assert.Equal(t, chess.DrawInsufficientMaterial, game.Result())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestResult_StalemateDetected validates US-05 / AC-05-04.
// Gherkin: "Library consumer detects stalemate when there are no legal moves and no check"
func TestResult_StalemateDetected(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(StalemateFEN)
	// require.NoError(t, err)
	// assert.Equal(t, chess.Stalemate, game.Result())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Castling ─────────────────────────────────────────────────────────────────

// TestCastling_WhiteKingsideApplied validates US-06 / AC-06-01.
// Gherkin: "Library consumer applies kingside castling for White"
func TestCastling_WhiteKingsideApplied(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(WhiteKingsideFEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e1g1")
	// newGame, err := game.Apply(move)
	// require.NoError(t, err)
	// assert.Equal(t, chess.WhiteKing, newGame.State.Board[chess.G1], "White king must be on g1")
	// assert.Equal(t, chess.WhiteRook, newGame.State.Board[chess.F1], "White rook must be on f1")
	// assert.False(t, newGame.State.CastlingRights&chess.CastleWhiteKingside != 0,
	//     "White kingside castling right must be removed")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestCastling_WhiteQueensideApplied validates AC-06-01 for queenside.
// Gherkin: "Library consumer applies queenside castling for White"
func TestCastling_WhiteQueensideApplied(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(WhiteQueensideFEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e1c1")
	// newGame, err := game.Apply(move)
	// require.NoError(t, err)
	// assert.Equal(t, chess.WhiteKing, newGame.State.Board[chess.C1])
	// assert.Equal(t, chess.WhiteRook, newGame.State.Board[chess.D1])

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestCastling_ThroughAttackedSquareRejected validates US-06 / AC-06-02.
// Gherkin: "Library consumer receives an error when castling through an attacked square"
func TestCastling_ThroughAttackedSquareRejected(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(CastlingThroughCheckFEN)
	// require.NoError(t, err)
	// illegalCastle := chess.Move{From: chess.E1, To: chess.G1}
	// _, err = game.Apply(illegalCastle)
	// require.Error(t, err)
	// assert.True(t, errors.Is(err, chess.ErrIllegalMove))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestCastling_WhileInCheckRejected validates US-06 / AC-06-03.
// Gherkin: "Library consumer receives an error when trying to castle while in check"
func TestCastling_WhileInCheckRejected(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// Position: White in check, castling rights still set.
	// inCheckWithCastlingFEN := "r3k2r/8/8/8/4r3/8/8/R3K2R w KQkq - 0 1"
	// game, err := chess.NewGameFromFEN(inCheckWithCastlingFEN)
	// require.NoError(t, err)
	// require.True(t, game.InCheck())
	// castle := chess.Move{From: chess.E1, To: chess.G1}
	// _, err = game.Apply(castle)
	// require.Error(t, err)
	// assert.True(t, errors.Is(err, chess.ErrIllegalMove))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── En Passant ───────────────────────────────────────────────────────────────

// TestEnPassant_CaptureRemovesCapturedPawn validates US-07 / AC-07-01.
// Gherkin: "Library consumer captures en passant and the captured pawn is removed"
func TestEnPassant_CaptureRemovesCapturedPawn(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// Position: Black just moved pawn to d5, White pawn on e5 can capture en passant on d6.
	epAvailableFEN := "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	_ = epAvailableFEN

	// game, err := chess.NewGameFromFEN(epAvailableFEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e5d6")
	// newGame, err := game.Apply(move)
	// require.NoError(t, err)
	// assert.Equal(t, chess.WhitePawn, newGame.State.Board[chess.D6], "White pawn must be on d6")
	// assert.Equal(t, chess.NoPiece, newGame.State.Board[chess.D5], "captured Black pawn must be removed")
	// assert.Equal(t, chess.NoSquare, newGame.State.EnPassantSq, "en passant square must be cleared")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestEnPassant_PinnedCaptureAbsentFromLegalMoves validates US-07 / AC-07-02.
// Gherkin: "Library consumer sees en passant absent from legal moves when the capture would expose the king"
func TestEnPassant_PinnedCaptureAbsentFromLegalMoves(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// A well-known en-passant pin position (Rook on e1, Pawn on e5, enemy pawn on d5, ep on d6,
	// but taking en passant exposes the king to the rook on a5).
	// epPinFEN := "8/8/8/r2pP2K/8/8/8/8 w - d6 0 1"
	// game, err := chess.NewGameFromFEN(epPinFEN)
	// require.NoError(t, err)
	// epMove := chess.Move{From: chess.E5, To: chess.D6}
	// for _, m := range game.LegalMoves() {
	//     if m == epMove {
	//         t.Fatal("en passant capture must be absent when it exposes the king")
	//     }
	// }

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Pawn Promotion ───────────────────────────────────────────────────────────

// TestPromotion_PawnBecomesQueen validates US-08 / AC-08-01.
// Gherkin: "Library consumer promotes a White pawn to a queen"
func TestPromotion_PawnBecomesQueen(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(PawnOnE7FEN)
	// require.NoError(t, err)
	// move := mustParseUCI(t, game, "e7e8q")
	// newGame, err := game.Apply(move)
	// require.NoError(t, err)
	// assert.Equal(t, chess.WhiteQueen, newGame.State.Board[chess.E8])
	// assert.Equal(t, chess.NoPiece, newGame.State.Board[chess.E7])

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPromotion_FourMovesGeneratedPerSquare validates US-08 / AC-08-02.
// Gherkin: "Library consumer sees four promotion moves for every reachable promotion square"
func TestPromotion_FourMovesGeneratedPerSquare(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(PawnOnE7FEN)
	// require.NoError(t, err)
	// moves := game.LegalMoves()
	// promoMoves := make([]chess.Move, 0)
	// for _, m := range moves {
	//     if m.IsPromotion() {
	//         promoMoves = append(promoMoves, m)
	//     }
	// }
	// assert.Len(t, promoMoves, 4, "expected 4 promotion moves for pawn on e7")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPromotion_MissingPieceReturnsError validates US-08 / AC-08-03.
// Gherkin: "Library consumer receives an error when applying a promotion move without specifying the piece"
func TestPromotion_MissingPieceReturnsError(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, err := chess.NewGameFromFEN(PawnOnE7FEN)
	// require.NoError(t, err)
	// noPromoMove := chess.Move{From: chess.E7, To: chess.E8, Promotion: chess.NoPiece}
	// _, err = game.Apply(noPromoMove)
	// require.Error(t, err)
	// assert.True(t, errors.Is(err, chess.ErrIllegalMove))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Notation Export ──────────────────────────────────────────────────────────

// TestUCIString_PawnMove validates US-09 / AC-09-01.
// Gherkin: "Library consumer receives correct UCI notation for a regular pawn move"
func TestUCIString_PawnMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// move := mustParseUCI(t, game, "e2e4")
	// assert.Equal(t, "e2e4", move.UCIString())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestUCIString_PromotionMove validates US-09 / AC-09-01.
// Gherkin: "Library consumer receives correct UCI notation for a promotion move"
func TestUCIString_PromotionMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(PawnOnE7FEN)
	// move := mustParseUCI(t, game, "e7e8q")
	// assert.Equal(t, "e7e8q", move.UCIString())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestUCIString_KingsideCastling validates US-09 / AC-09-01.
// Gherkin: "Library consumer receives correct UCI notation for kingside castling"
func TestUCIString_KingsideCastling(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(WhiteKingsideFEN)
	// move := mustParseUCI(t, game, "e1g1")
	// assert.Equal(t, "e1g1", move.UCIString())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_PawnMove validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for a pawn move"
func TestSANString_PawnMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// move := mustParseUCI(t, game, "e2e4")
	// assert.Equal(t, "e4", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_KnightMove validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for a knight move"
func TestSANString_KnightMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// move := mustParseUCI(t, game, "g1f3")
	// assert.Equal(t, "Nf3", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_KingsideCastling validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for kingside castling"
func TestSANString_KingsideCastling(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(WhiteKingsideFEN)
	// move := mustParseUCI(t, game, "e1g1")
	// assert.Equal(t, "O-O", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_QueensideCastling validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for queenside castling"
func TestSANString_QueensideCastling(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(WhiteQueensideFEN)
	// move := mustParseUCI(t, game, "e1c1")
	// assert.Equal(t, "O-O-O", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_PromotionMove validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for a promotion move"
func TestSANString_PromotionMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(PawnOnE7FEN)
	// move := mustParseUCI(t, game, "e7e8q")
	// assert.Equal(t, "e8=Q", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSANString_CheckmateMove validates US-10 / AC-09-02.
// Gherkin: "Library consumer receives correct SAN notation for a checkmate move"
func TestSANString_CheckmateMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// preFoolsMateFEN = one move before Fool's Mate
	preFoolsMateFEN := "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq g3 0 2"
	_ = preFoolsMateFEN

	// game, _ := chess.NewGameFromFEN(preFoolsMateFEN)
	// move := mustParseUCI(t, game, "d8h4")
	// assert.Equal(t, "Qh4#", move.SANString(game))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── PGN Export ───────────────────────────────────────────────────────────────

// TestPGNExport_CompleteGameContainsRequiredSections validates US-11 / AC-10-01.
// Gherkin: "Library consumer exports a complete game as a valid PGN string"
func TestPGNExport_CompleteGameContainsRequiredSections(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// Build Fool's Mate move by move.
	// moves := []string{"f2f3", "e7e5", "g2g4", "d8h4"}
	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// for _, uci := range moves {
	//     m := mustParseUCI(t, game, uci)
	//     game, _ = game.Apply(m)
	// }
	// pgn := game.ToPGN()
	// for _, tag := range []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"} {
	//     assert.Contains(t, pgn, "["+tag, "PGN missing tag: %s", tag)
	// }
	// assert.Contains(t, pgn, "0-1", "PGN must contain result token 0-1 for Black win")
	// assert.Regexp(t, `1\. `, pgn, "PGN must contain move number 1")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPGNExport_InProgressGameHasStarResult validates US-11 / AC-10-01.
// Gherkin: "Library consumer exports an in-progress game and the PGN has a star result token"
func TestPGNExport_InProgressGameHasStarResult(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// for _, uci := range []string{"e2e4", "e7e5", "g1f3"} {
	//     m := mustParseUCI(t, game, uci)
	//     game, _ = game.Apply(m)
	// }
	// pgn := game.ToPGN()
	// assert.Contains(t, pgn, "*", "in-progress PGN must end with *")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Perft Validation ─────────────────────────────────────────────────────────

// TestPerft_StartingPositionDepth1 validates US-12 / AC-11-01.
// Gherkin: "Move generator produces exactly 20 nodes at depth 1 from the starting position"
func TestPerft_StartingPositionDepth1(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// assert.Equal(t, 20, perft(game, 1))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_StartingPositionDepth2 validates US-12 / AC-11-01.
// Gherkin: "Move generator produces exactly 400 nodes at depth 2 from the starting position"
func TestPerft_StartingPositionDepth2(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// assert.Equal(t, 400, perft(game, 2))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_StartingPositionDepth3 validates US-12 / AC-11-01.
// Gherkin: "Move generator produces exactly 8902 nodes at depth 3 from the starting position"
func TestPerft_StartingPositionDepth3(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// assert.Equal(t, 8902, perft(game, 3))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_StartingPositionDepth4 validates US-12 / AC-11-01.
// Gherkin: "Move generator produces exactly 197281 nodes at depth 4 from the starting position"
func TestPerft_StartingPositionDepth4(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// assert.Equal(t, 197281, perft(game, 4))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_StartingPositionDepth5 validates US-12 / AC-11-01.
// Gherkin: "Move generator produces exactly 4865609 nodes at depth 5 from the starting position"
// Tagged @slow in the feature file — this test is gated by -run TestPerft_StartingPositionDepth5
func TestPerft_StartingPositionDepth5(t *testing.T) {
	skipUnimplemented(t)
	if testing.Short() {
		t.Skip("skipping depth-5 perft in short mode")
	}
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// assert.Equal(t, 4865609, perft(game, 5))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_KiwipeteDepth1 validates US-12 / AC-11-02.
// Gherkin: "Move generator produces correct node counts from the Kiwipete position at depth 1"
func TestPerft_KiwipeteDepth1(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(KiwipeteFEN)
	// assert.Equal(t, 48, perft(game, 1))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_KiwipeteDepth2 validates US-12 / AC-11-02.
// Gherkin: "Move generator produces correct node counts from the Kiwipete position at depth 2"
func TestPerft_KiwipeteDepth2(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(KiwipeteFEN)
	// assert.Equal(t, 2039, perft(game, 2))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_KiwipeteDepth3 validates US-12 / AC-11-02.
// Gherkin: "Move generator produces correct node counts from the Kiwipete position at depth 3"
func TestPerft_KiwipeteDepth3(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(KiwipeteFEN)
	// assert.Equal(t, 97862, perft(game, 3))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestPerft_KiwipeteDepth4 validates US-12 / AC-11-02.
// Gherkin: "Move generator produces correct node counts from the Kiwipete position at depth 4"
func TestPerft_KiwipeteDepth4(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess")

	// game, _ := chess.NewGameFromFEN(KiwipeteFEN)
	// assert.Equal(t, 4085603, perft(game, 4))

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// perft counts nodes at the given depth using the chess package LegalMoves and Apply.
// Uncomment and use once internal/chess is implemented.
//
//nolint:unused
func perft_helper(t *testing.T, depth int, fen string) int {
	t.Helper()
	_ = fen
	_ = depth
	// game, err := chess.NewGameFromFEN(fen)
	// require.NoError(t, err)
	// return perft(game, depth)
	return 0
}

// mustParseUCI finds the move with the given UCI string in game.LegalMoves().
// It fails the test if the move is not found.
//
//nolint:unused
func mustParseUCI_helper(t *testing.T, uci string) {
	t.Helper()
	_ = uci
	// for _, m := range game.LegalMoves() {
	//     if m.UCIString() == uci {
	//         return m
	//     }
	// }
	// t.Fatalf("move %q not found in legal moves", uci)
	// panic("unreachable")
}

// containsSubstring is a helper used throughout the step files.
func containsSubstring(s, sub string) bool {
	return strings.Contains(s, sub)
}
