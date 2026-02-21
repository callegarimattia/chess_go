// tui_steps_test.go — Executable specifications for the TUI game loop.
//
// Mirrors: milestone-3-tui.feature
// Driving port: tui.NewGame(r io.Reader, w io.Writer, engineFn EngineFunc) tui.Game
//   called via tui.Game.Run()
//
// The engine is injected as a stub (tui.EngineFunc) for determinism.
// No real subprocess is launched. No os.Stdin or os.Stdout is used.
//
// CM-A compliance: all calls go through tui.NewGame public API.
// CM-B compliance: test names use domain language exclusively.

package acceptance_test

import (
	"bytes"
	"strings"
	"testing"
	// Production packages — uncomment once implemented.
	// chess "chess_go/internal/chess"
	// engine "chess_go/internal/engine"
	// "chess_go/internal/tui"
)

// ─── Stub Engine ──────────────────────────────────────────────────────────────

// stubEngine returns a fixed UCI move string regardless of position.
// It satisfies the tui.EngineFunc signature once the tui package is implemented.
//
//nolint:unused
type stubEngineFunc struct {
	response string // UCI string of the move to always return, e.g. "e7e5"
	called   int    // number of times the engine was called
}

// call is a placeholder for tui.EngineFunc invocation.
// Uncomment and adapt to tui.EngineFunc signature once the tui package exists.
//
//nolint:unused
func (s *stubEngineFunc) call() {
	s.called++
	// return chess.Move matching s.response
}

// ─── Board Render ─────────────────────────────────────────────────────────────

// TestTUIRender_StartingBoardShowsAllRequiredElements validates US-24 / AC-15-01.
// Gherkin: "TUI player launches the game and sees the starting board"
func TestTUIRender_StartingBoardShowsAllRequiredElements(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer
	_ = input

	// Provide EOF immediately so the game loop exits after rendering.
	// input stays empty.

	// stub := func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     return chess.Move{} // never called — input is empty
	// }
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "White to move", "must show side-to-move indicator")
	// assertContainsCoordinates(t, rendered)
	// assertContainsPieces(t, rendered)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIRender_AllStartingPiecesVisible validates US-24 / AC-15-01.
// Gherkin: "TUI player sees all starting pieces in their correct positions"
func TestTUIRender_AllStartingPiecesVisible(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// game := tui.NewGame(&input, &output, noOpEngine)
	// game.Run()
	rendered := output.String()
	_ = rendered

	// The render must contain piece symbols for all 32 starting pieces.
	// At minimum, king symbols for both sides must appear.
	// assert.Contains(t, rendered, "K", "White king symbol must be present")
	// assert.Contains(t, rendered, "k", "Black king symbol must be present (or Unicode equivalent)")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIRender_BlackToMoveIndicatorAfterFirstMove validates US-24 / AC-15-01.
// Gherkin: "TUI player sees Black to move indicator after the first move"
func TestTUIRender_BlackToMoveIndicatorAfterFirstMove(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess", "internal/engine")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Provide one valid move then EOF.
	input.WriteString("e2e4\n")

	// stub := func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     // Return a legal Black move so the loop advances.
	//     m, _ := findMove(g, "e7e5")
	//     return m
	// }
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Black to move",
	//     "after White's first move, board must show Black to move")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Legal Move Acceptance ────────────────────────────────────────────────────

// TestTUIInput_LegalMoveUpdatesBoard validates US-25 / AC-15-02.
// Gherkin: "TUI player enters a legal move and sees the updated board"
func TestTUIInput_LegalMoveUpdatesBoard(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess", "internal/engine")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	input.WriteString("e2e4\n")

	// stub := alwaysRespondsWithMove("e7e5")
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Engine thinking",
	//     "TUI must show thinking indicator while engine computes")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIInput_LegalMoveAndEngineResponseBothReflected validates US-25, US-26 / AC-15-02.
// Gherkin: "TUI player enters a legal move and the engine responds with a move"
func TestTUIInput_LegalMoveAndEngineResponseBothReflected(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess", "internal/engine")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	input.WriteString("e2e4\n")

	// stub := alwaysRespondsWithMove("e7e5")
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// The board re-render after engine response must show both moves.
	// assertBoardContainsPieceOnSquare(t, rendered, "e4") // White pawn
	// assertBoardContainsPieceOnSquare(t, rendered, "e5") // Black pawn (after engine)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Illegal Move Rejection ───────────────────────────────────────────────────

// TestTUIInput_IllegalMoveShowsErrorAndRepromptsUser validates US-25 / AC-15-03.
// Gherkin: "TUI player enters an illegal move and sees a clear error message"
func TestTUIInput_IllegalMoveShowsErrorAndRepromptsUser(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// First send an illegal move, then a quit/EOF to end the loop.
	input.WriteString("e2e5\n")

	// game := tui.NewGame(&input, &output, noOpEngine)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Illegal move",
	//     "error message must contain 'Illegal move'")
	// assert.Contains(t, rendered, "e2e5",
	//     "error message must echo the attempted move string")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIInput_InvalidFormatShowsFormatError validates US-25 / AC-15-03.
// Gherkin: "TUI player enters a move in the wrong format and sees a format error"
func TestTUIInput_InvalidFormatShowsFormatError(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	input.WriteString("hello\n")

	// game := tui.NewGame(&input, &output, noOpEngine)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.True(t,
	//     strings.Contains(rendered, "invalid") || strings.Contains(rendered, "format"),
	//     "error message must indicate invalid format; got: %q", rendered)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIInput_EmptyInputShowsPromptAgain validates AC-15-03 edge case.
// Gherkin: "TUI player enters an empty input and the prompt is shown again"
func TestTUIInput_EmptyInputShowsPromptAgain(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Send empty line then EOF.
	input.WriteString("\n")

	// game := tui.NewGame(&input, &output, noOpEngine)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// Count move prompts: at minimum two (original + re-prompt after empty input).
	// promptCount := strings.Count(rendered, "Your move:")
	// assert.GreaterOrEqual(t, promptCount, 2)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIInput_MovingOpponentPieceIsRejected validates AC-15-03.
// Gherkin: "TUI player tries to move a piece that belongs to the opponent"
func TestTUIInput_MovingOpponentPieceIsRejected(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// White to move; e7e5 is a Black move.
	input.WriteString("e7e5\n")

	// game := tui.NewGame(&input, &output, noOpEngine)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Illegal move")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Check Notification ───────────────────────────────────────────────────────

// TestTUIStatus_CheckNotificationDisplayed validates US-27 / AC-15-04.
// Gherkin: "TUI player sees a check notification when a move delivers check"
func TestTUIStatus_CheckNotificationDisplayed(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess", "internal/engine")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Position: White plays Qh5+ delivering check. Engine stub returns a checking move.
	// checkingPosFEN := "rnbqkb1r/pppp1ppp/5n2/4p1Q1/4P3/8/PPPP1PPP/RNB1KBNR b KQkq - 3 3"
	// input stays empty — engine delivers check, game detects it.

	// stub := func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     m, _ := findMove(g, "d1h5") // Qh5+
	//     return m
	// }
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Check!")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Game Result Display ─────────────────────────────────────────────────────

// TestTUIResult_CheckmateDisplayedAndNoFurtherPrompt validates US-27 / AC-15-05.
// Gherkin: "TUI player sees a checkmate result message and no further move prompt"
func TestTUIResult_CheckmateDisplayedAndNoFurtherPrompt(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess", "internal/engine")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Position: engine's first move is checkmate.
	// preMateFEN := "k7/8/1K6/8/8/8/8/R7 b - - 0 1"
	// Black to move with no legal moves → engine is not called in this direction.
	// Actually: use a position where engine delivers mate as White.
	// Simplest: start game with empty input, engine stub always returns Ra8#.

	// stub := func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     m, _ := findMove(g, "a1a8") // Ra8# in the Lucena-style position
	//     return m
	// }
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Checkmate", "result message must contain 'Checkmate'")
	// promptsAfterResult := countPromptsAfterSubstring(rendered, "Checkmate", "Your move:")
	// assert.Equal(t, 0, promptsAfterResult, "no move prompt must appear after checkmate")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIResult_StalemateMessageDisplayed validates US-27 / AC-15-05.
// Gherkin: "TUI player sees a stalemate result message when the game is drawn by stalemate"
func TestTUIResult_StalemateMessageDisplayed(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Load stalemate position directly (no engine needed — no legal moves).
	// game := tui.NewGameFromFEN(&input, &output, noOpEngine, StalemateFEN)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Stalemate")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIResult_DrawByFiftyMoveRuleDisplayed validates US-27 / AC-15-05.
// Gherkin: "TUI player sees a draw-by-fifty-move-rule message"
func TestTUIResult_DrawByFiftyMoveRuleDisplayed(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// FEN with half-move clock at 100.
	fiftyMoveFEN := "8/8/8/8/8/8/8/K6k w - - 100 101"
	_ = fiftyMoveFEN

	// game := tui.NewGameFromFEN(&input, &output, noOpEngine, fiftyMoveFEN)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.True(t,
	//     strings.Contains(rendered, "fifty") || strings.Contains(rendered, "draw"),
	//     "must mention fifty-move rule")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIResult_DrawByRepetitionDisplayed validates US-27 / AC-15-05.
// Gherkin: "TUI player sees a draw-by-threefold-repetition message"
func TestTUIResult_DrawByRepetitionDisplayed(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var output bytes.Buffer

	// Simulate a game reaching threefold repetition via the TUI loop.
	// Moves: Ng1f3 Ng8f6 Nf3g1 Nf6g8 Ng1f3 Ng8f6 Nf3g1 Nf6g8 (same position 3 times).
	repetitionMoves := []string{"g1f3", "g8f6", "f3g1", "f6g8", "g1f3", "g8f6", "f3g1", "f6g8"}
	var inputBuf bytes.Buffer
	for _, m := range repetitionMoves {
		inputBuf.WriteString(m + "\n")
	}
	_ = inputBuf
	_ = output

	// stub := alternatingMoveStub(t, repetitionMoves)
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.True(t,
	//     strings.Contains(rendered, "repetition") || strings.Contains(rendered, "draw"),
	//     "must mention threefold repetition")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── PGN Export ───────────────────────────────────────────────────────────────

// TestTUIPGN_SavePromptAppearsAfterGameEnds validates US-28 / AC-15-05.
// Gherkin: "TUI player is offered a save prompt when the game ends"
func TestTUIPGN_SavePromptAppearsAfterGameEnds(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	// Position where engine immediately delivers checkmate.
	// After checkmate, TUI must offer PGN save before exiting.
	// Send 'n' to decline the save.
	input.WriteString("n\n")

	// stub := mateInOneEngineStub(t)
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	rendered := output.String()
	_ = rendered

	// assert.Contains(t, rendered, "Save game", "save prompt must appear after game ends")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIPGN_SaveWritesPGNToFile validates US-28 / AC-15-05.
// Gherkin: "TUI player saves the PGN and the file is written to disk"
func TestTUIPGN_SaveWritesPGNToFile(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	outFile := t.TempDir() + "/game.pgn"
	// Answer "y" then provide file path.
	input.WriteString("y\n")
	input.WriteString(outFile + "\n")

	// stub := mateInOneEngineStub(t)
	// game := tui.NewGame(&input, &output, stub)
	// game.Run()

	// _, err := os.Stat(outFile)
	// require.NoError(t, err, "PGN file must be created")
	// content, err := os.ReadFile(outFile)
	// require.NoError(t, err)
	// pgn := string(content)
	// for _, tag := range []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"} {
	//     assert.Contains(t, pgn, "["+tag, "PGN missing tag: %s", tag)
	// }

	_ = outFile
	_ = output

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTUIPGN_DeclineSaveExitsCleanly validates US-28 / AC-15-05.
// Gherkin: "TUI player declines to save the PGN and the game exits cleanly"
func TestTUIPGN_DeclineSaveExitsCleanly(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/tui", "internal/chess")

	var input bytes.Buffer
	_ = input
	var output bytes.Buffer

	input.WriteString("n\n")

	// stub := mateInOneEngineStub(t)
	// game := tui.NewGame(&input, &output, stub)
	// exitCode := game.Run()
	// assert.Equal(t, 0, exitCode)
	// assert.NotContains(t, output.String(), "panic")

	_ = output

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// assertContainsCoordinates checks that rank and file labels appear in the rendered output.
//
//nolint:unused
func assertContainsCoordinates(t *testing.T, rendered string) {
	t.Helper()
	for _, file := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		if !strings.Contains(rendered, file) {
			t.Errorf("rendered board missing file label %q", file)
		}
	}
	for _, rank := range []string{"1", "2", "3", "4", "5", "6", "7", "8"} {
		if !strings.Contains(rendered, rank) {
			t.Errorf("rendered board missing rank label %q", rank)
		}
	}
}

// assertContainsPieces checks that at least some piece symbols appear in the output.
//
//nolint:unused
func assertContainsPieces(t *testing.T, rendered string) {
	t.Helper()
	// ASCII or Unicode — at minimum the king symbol must appear.
	hasPieces := strings.ContainsAny(rendered, "KQRBNPkqrbnp♔♕♖♗♘♙♚♛♜♝♞♟")
	if !hasPieces {
		t.Errorf("rendered board contains no piece symbols")
	}
}

// noOpEngine is a stub engine function that panics if called.
// Use it when the test does not expect the engine to be invoked.
//
//nolint:unused
func noOpEngine_helper(t *testing.T) {
	t.Helper()
	// return func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     t.Fatal("engine must not be called in this scenario")
	//     return chess.Move{}
	// }
}
