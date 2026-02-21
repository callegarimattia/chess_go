// engine_steps_test.go — Executable specifications for the engine and UCI protocol.
//
// Mirrors: milestone-2-engine.feature
// Driving ports:
//   - engine.Search(g chess.Game, tc engine.TimeControl, info io.Writer) engine.SearchResult
//   - engine.UCIHandler.Run(r io.Reader, w io.Writer)
//   - chess-go binary via os/exec (for UCI subprocess tests)
//
// CM-A compliance: calls engine.Search and engine.UCIHandler only through public API.
//   UCI subprocess tests invoke the binary, not internal packages.
// CM-B compliance: test names use domain language exclusively.

package acceptance_test

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
	// Production packages — uncomment once implemented.
	// chess "chess_go/internal/chess"
	// engine "chess_go/internal/engine"
)

// ─── Random Move Engine ───────────────────────────────────────────────────────

// TestRandomEngine_ReturnsLegalMoveFromStartingPosition validates US-13.
// Gherkin: "Engine developer gets a legal move from the random move selector"
func TestRandomEngine_ReturnsLegalMoveFromStartingPosition(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 10 * time.Millisecond}
	// var info bytes.Buffer
	// start := time.Now()
	// result := engine.Search(game, tc, &info)
	// elapsed := time.Since(start)
	//
	// assertWithinDuration(t, 50*time.Millisecond, elapsed, "random engine must return within 50ms")
	//
	// legal := game.LegalMoves()
	// found := false
	// for _, m := range legal {
	//     if m.UCIString() == result.BestMove.UCIString() {
	//         found = true
	//         break
	//     }
	// }
	// assert.True(t, found, "returned move %q must be in LegalMoves()", result.BestMove.UCIString())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestRandomEngine_ReturnsLegalMoveFromKiwipete validates US-13.
// Gherkin: "Engine developer gets a legal move from any non-terminal position"
func TestRandomEngine_ReturnsLegalMoveFromKiwipete(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(KiwipeteFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 10 * time.Millisecond}
	// result := engine.Search(game, tc, io.Discard)
	// legal := game.LegalMoves()
	// found := false
	// for _, m := range legal {
	//     if m.UCIString() == result.BestMove.UCIString() {
	//         found = true; break
	//     }
	// }
	// assert.True(t, found)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Alpha-Beta Search ────────────────────────────────────────────────────────

// TestSearch_LegalBestmoveWithinTimeLimitFromStart validates US-14 / AC-12-01.
// Gherkin: "Engine developer receives a legal bestmove within the time limit from the starting position"
func TestSearch_LegalBestmoveWithinTimeLimitFromStart(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 1000 * time.Millisecond}
	// var info bytes.Buffer
	// start := time.Now()
	// result := engine.Search(game, tc, &info)
	// elapsed := time.Since(start)
	//
	// assertWithinDuration(t, 1050*time.Millisecond, elapsed)
	//
	// legal := game.LegalMoves()
	// assertMoveIsLegal(t, result.BestMove, legal)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSearch_EmitsInfoLinesDuringSearch validates US-14 / AC-12-01.
// Gherkin: "Engine developer sees info lines emitted during search"
func TestSearch_EmitsInfoLinesDuringSearch(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 500 * time.Millisecond}
	// var info bytes.Buffer
	// _ = engine.Search(game, tc, &info)
	//
	// lines := strings.Split(strings.TrimSpace(info.String()), "\n")
	// hasInfoLine := false
	// for _, line := range lines {
	//     if strings.HasPrefix(line, "info depth") {
	//         hasInfoLine = true
	//         assert.Contains(t, line, "score")
	//         assert.Contains(t, line, "nodes")
	//         assert.Contains(t, line, "nps")
	//         assert.Contains(t, line, "pv")
	//     }
	// }
	// assert.True(t, hasInfoLine, "at least one info depth line must be emitted")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSearch_ReachesDepth3In100ms validates US-14 / AC-12-02.
// Gherkin: "Engine developer sees the search reach at least depth 3 in 100 milliseconds"
func TestSearch_ReachesDepth3In100ms(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 100 * time.Millisecond}
	// var info bytes.Buffer
	// result := engine.Search(game, tc, &info)
	//
	// assert.GreaterOrEqual(t, result.Depth, 3, "engine must reach at least depth 3 in 100ms")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSearch_FindsMateInOne validates US-14 / AC-12-03.
// Gherkin: "Engine developer sees the engine find a forced mate in one"
func TestSearch_FindsMateInOne(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(MateIn1FEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 100 * time.Millisecond}
	// result := engine.Search(game, tc, io.Discard)
	// applied, err := game.Apply(result.BestMove)
	// require.NoError(t, err)
	// assert.Equal(t, chess.WhiteWins, applied.Result(), "engine must find the mating move")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSearch_FindsFoolsMateMoveAsBlack validates AC-12-03 with Fool's Mate.
// Gherkin: "Engine developer sees the engine find a forced mate in one in Fool's Mate setup"
func TestSearch_FindsFoolsMateMoveAsBlack(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	preFoolsMate := "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq g3 0 2"
	_ = preFoolsMate

	// game, err := chess.NewGameFromFEN(preFoolsMate)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 100 * time.Millisecond}
	// result := engine.Search(game, tc, io.Discard)
	// assert.Equal(t, "d8h4", result.BestMove.UCIString(), "engine must find Qh4# (Fool's Mate)")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Time Management ─────────────────────────────────────────────────────────

// TestTimeManagement_BestmoveWithinGracePeriod validates US-17 / AC-13-01.
// Gherkin: "Engine developer sees the bestmove returned within the movetime grace period"
func TestTimeManagement_BestmoveWithinGracePeriod(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	movetime := 500 * time.Millisecond
	grace := 50 * time.Millisecond

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: movetime}
	// start := time.Now()
	// _ = engine.Search(game, tc, io.Discard)
	// elapsed := time.Since(start)
	// assertWithinDuration(t, movetime+grace, elapsed)

	t.Log("time management: movetime=", movetime, " grace=", grace)
	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestTimeManagement_VeryShortMovetime validates US-17 / AC-13-01.
// Gherkin: "Engine developer sees the engine respect a very short movetime"
func TestTimeManagement_VeryShortMovetime(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/chess", "internal/engine")

	// game, err := chess.NewGameFromFEN(StartingFEN)
	// require.NoError(t, err)
	// tc := engine.TimeControl{MoveTime: 50 * time.Millisecond}
	// start := time.Now()
	// result := engine.Search(game, tc, io.Discard)
	// elapsed := time.Since(start)
	// assertWithinDuration(t, 100*time.Millisecond, elapsed)
	// assertMoveIsLegal(t, result.BestMove, game.LegalMoves())

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── UCI Handshake ─────────────────────────────────────────────────────────────

// TestUCI_HandshakeReturnsRequiredLines validates US-20 / AC-14-01.
// Gherkin: "Engine developer sends uci and receives the required identification lines"
func TestUCI_HandshakeReturnsRequiredLines(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)

	start := time.Now()
	sess.send("uci")

	lines, found := sess.readUntil("uciok", 100*time.Millisecond)
	elapsed := time.Since(start)

	if !found {
		t.Fatalf("uciok not received within 100ms; got lines: %v", lines)
	}
	assertWithinDuration(t, 100*time.Millisecond, elapsed)

	combined := strings.Join(lines, "\n")
	if !containsSubstring(combined, "id name chess-go") {
		t.Errorf("expected 'id name chess-go' in UCI response; got:\n%s", combined)
	}
	if !containsSubstring(combined, "id author") {
		t.Errorf("expected 'id author' line in UCI response; got:\n%s", combined)
	}
}

// TestUCI_IsreadyReturnsReadyok validates US-20 / AC-14-02.
// Gherkin: "Engine developer sends isready after uci and receives readyok"
func TestUCI_IsreadyReturnsReadyok(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)

	start := time.Now()
	sess.send("isready")
	lines, found := sess.readUntil("readyok", 100*time.Millisecond)
	elapsed := time.Since(start)

	if !found {
		t.Fatalf("readyok not received; got: %v", lines)
	}
	assertWithinDuration(t, 100*time.Millisecond, elapsed)
}

// TestUCI_UCINewGameResetsState validates US-20 / AC-14 (ucinewgame).
// Gherkin: "Engine developer resets state with ucinewgame between two searches"
func TestUCI_UCINewGameResetsState(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	// First search.
	sess.send("position startpos")
	sess.send("go movetime 100")
	sess.readUntil("bestmove", 300*time.Millisecond)

	// Reset.
	sess.send("ucinewgame")
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	// Second search — must return a valid bestmove.
	sess.send("position startpos")
	sess.send("go movetime 100")
	lines, found := sess.readUntil("bestmove", 300*time.Millisecond)
	if !found {
		t.Fatalf("bestmove not received after ucinewgame; got: %v", lines)
	}
	combined := strings.Join(lines, "\n")
	if !containsSubstring(combined, "bestmove") {
		t.Errorf("expected bestmove in output; got:\n%s", combined)
	}
}

// ─── Position Command ─────────────────────────────────────────────────────────

// TestUCI_PositionStartposMoves validates US-21 / AC-14-03.
// Gherkin: "Engine developer sets up a position from the starting position with moves applied"
func TestUCI_PositionStartposMoves(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	sess.send("position startpos moves e2e4 e7e5")
	sess.send("go movetime 500")
	lines, found := sess.readUntil("bestmove", 600*time.Millisecond)
	if !found {
		t.Fatalf("bestmove not received; got: %v", lines)
	}

	// Extract bestmove token and verify it's legal in the position after e2e4 e7e5.
	// game, _ := chess.NewGameFromFEN(StartingFEN)
	// game, _ = game.Apply(mustParseUCI(t, game, "e2e4"))
	// game, _ = game.Apply(mustParseUCI(t, game, "e7e5"))
	// bestmoveUCI := extractBestmove(lines)
	// assertMoveIsLegal(t, chess.Move{}, game.LegalMoves()) // replace with parsed move

	_ = lines
}

// TestUCI_PositionFEN validates US-21 / AC-14-03 (FEN variant).
// Gherkin: "Engine developer sets up a position directly from a FEN string"
func TestUCI_PositionFEN(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	sess.send(fmt.Sprintf("position fen %s", AfterE2E4FEN))
	sess.send("go movetime 500")
	lines, found := sess.readUntil("bestmove", 600*time.Millisecond)
	if !found {
		t.Fatalf("bestmove not received after position fen; got: %v", lines)
	}
	_ = lines
}

// ─── Stop Command ─────────────────────────────────────────────────────────────

// TestUCI_StopCommandYieldsBestmoveWithin100ms validates US-23 / AC-14-04.
// Gherkin: "Engine developer sends stop during an active search and receives bestmove promptly"
func TestUCI_StopCommandYieldsBestmoveWithin100ms(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	sess.send("position startpos")
	sess.send("go infinite")
	time.Sleep(200 * time.Millisecond)

	start := time.Now()
	sess.send("stop")
	lines, found := sess.readUntil("bestmove", 200*time.Millisecond)
	elapsed := time.Since(start)

	if !found {
		t.Fatalf("bestmove not received after stop; got: %v", lines)
	}
	assertWithinDuration(t, 100*time.Millisecond, elapsed, "bestmove must arrive within 100ms of stop")
}

// TestUCI_QuitExitsCleanly validates US-23 / AC-14 (quit).
// Gherkin: "Engine developer sends quit and the process exits cleanly"
func TestUCI_QuitExitsCleanly(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)

	sess.send("quit")

	done := make(chan error, 1)
	go func() { done <- sess.cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("process exited with error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("process did not exit within 500ms of quit")
	}
}

// TestUCI_UnknownCommandIgnoredGracefully validates US-23 / AC-14-05.
// Gherkin: "Engine developer sends an unknown command and the engine does not crash"
func TestUCI_UnknownCommandIgnoredGracefully(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine", "cmd/chess-go")

	binPath := mustBuildBinary(t, "./cmd/chess-go")
	sess := newUCISession(t, binPath)
	sess.send("uci")
	sess.readUntil("uciok", 200*time.Millisecond)
	sess.send("isready")
	sess.readUntil("readyok", 100*time.Millisecond)

	// Send unknown command — should produce no output.
	sess.send("foo bar baz")
	time.Sleep(50 * time.Millisecond)

	// Engine still alive: isready returns readyok.
	sess.send("isready")
	lines, found := sess.readUntil("readyok", 100*time.Millisecond)
	if !found {
		t.Fatalf("engine did not respond to isready after unknown command; got: %v", lines)
	}
}

// ─── UCIHandler in-process (alternative to subprocess) ────────────────────────

// TestUCIHandler_HandshakeInProcess validates the UCI handler driving port without subprocess.
// This exercises engine.UCIHandler directly for faster feedback in the TDD loop.
// Gherkin: milestone-2-engine.feature — "UCI handshake" (in-process variant)
func TestUCIHandler_HandshakeInProcess(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine")

	var input bytes.Buffer
	var output bytes.Buffer

	input.WriteString("uci\n")
	input.WriteString("quit\n")

	// handler := engine.NewUCIHandler(engine.Search)
	// handler.Run(&input, &output)
	//
	// response := output.String()
	// assert.Contains(t, response, "id name chess-go")
	// assert.Contains(t, response, "id author")
	// assert.Contains(t, response, "uciok")

	_ = input
	_ = output

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestUCIHandler_PositionAndGoInProcess validates position+go in-process.
// Gherkin: "Engine developer sets up a position from the starting position with moves applied" (in-process)
func TestUCIHandler_PositionAndGoInProcess(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/engine")

	var input bytes.Buffer
	var output bytes.Buffer

	input.WriteString("uci\n")
	input.WriteString("isready\n")
	input.WriteString("position startpos moves e2e4 e7e5\n")
	input.WriteString("go movetime 200\n")
	input.WriteString("quit\n")

	// handler := engine.NewUCIHandler(engine.Search)
	// handler.Run(&input, &output)
	//
	// response := output.String()
	// assert.Contains(t, response, "bestmove", "in-process UCI handler must emit bestmove")
	// bestmoveUCI := extractBestmove(strings.Split(response, "\n"))
	// // verify it's legal in the position after e2e4 e7e5
	// _ = bestmoveUCI

	_ = bufio.NewReader(&input)
	_ = output

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// extractBestmove finds the UCI bestmove from a slice of output lines.
//
//nolint:unused
func extractBestmove(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "bestmove ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1]
			}
		}
	}
	return ""
}

// assertMoveIsLegal verifies that move appears in the legal move list.
// Placeholder — uncomment and adapt once chess.Move is available.
//
//nolint:unused
func assertMoveIsLegalHelper(t *testing.T, uci string) {
	t.Helper()
	_ = uci
}
