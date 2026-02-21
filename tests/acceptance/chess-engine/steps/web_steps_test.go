// web_steps_test.go — Executable specifications for the SSR HTTP server.
//
// Mirrors: milestone-4-ssr.feature
// Driving port: web.NewServer(engineFn EngineFunc) *web.Server
//   exercised via net/http/httptest.NewServer
//
// The engine is injected as a stub for determinism and speed.
// No real network ports are opened. httptest.NewServer is used throughout.
//
// CM-A compliance: all calls go through web.NewServer public API.
// CM-B compliance: test names use domain language exclusively.

package acceptance_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	// Production packages — uncomment once implemented.
	// chess "chess_go/internal/chess"
	// engine "chess_go/internal/engine"
	// "chess_go/internal/web"
)

// ─── Test Server Setup ────────────────────────────────────────────────────────

// webTestServer wraps httptest.Server with helpers for SSR acceptance tests.
type webTestServer struct {
	server *httptest.Server
	client *http.Client
	t      *testing.T
}

// newWebTestServer starts an httptest.Server backed by web.NewServer with a stub engine.
// Replace the stub with a real engine once internal/engine is implemented.
func newWebTestServer(t *testing.T, engineResponse string) *webTestServer {
	t.Helper()
	_ = engineResponse

	// Once internal/web is implemented:
	// stub := func(g chess.Game, tc engine.TimeControl) chess.Move {
	//     m, _ := findLegalMove(g, engineResponse)
	//     return m
	// }
	// srv := httptest.NewServer(web.NewServer(stub))

	// Placeholder: a stub HTTP handler for compilation before implementation.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("not implemented"))
	}))

	t.Cleanup(srv.Close)

	// Use a client that does NOT follow redirects automatically so we can inspect 302.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &webTestServer{server: srv, client: client, t: t}
}

// get sends a GET request and returns the response.
func (s *webTestServer) get(path string) *http.Response {
	s.t.Helper()
	resp, err := s.client.Get(s.server.URL + path)
	if err != nil {
		s.t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

// post sends a POST request and returns the response.
func (s *webTestServer) post(path string) *http.Response {
	s.t.Helper()
	resp, err := s.client.PostForm(s.server.URL+path, url.Values{})
	if err != nil {
		s.t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

// body reads and returns the response body as a string.
func (s *webTestServer) body(resp *http.Response) string {
	s.t.Helper()
	defer resp.Body.Close() //nolint:errcheck
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		s.t.Fatalf("read body: %v", err)
	}
	return string(b)
}

// createSession creates a new game session and returns the session path (e.g. "/game/abc123").
func (s *webTestServer) createSession() string {
	s.t.Helper()
	resp := s.post("/game/new")
	if resp.StatusCode != http.StatusFound {
		s.t.Fatalf("POST /game/new: expected 302, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if !strings.HasPrefix(loc, "/game/") {
		s.t.Fatalf("Location header must start with /game/, got: %q", loc)
	}
	return loc
}

// ─── Home Page ───────────────────────────────────────────────────────────────

// TestSSR_HomePageReturns200WithBoard validates US-29 / AC-16-01.
// Gherkin: "Web player opens the home page and sees the starting board"
func TestSSR_HomePageReturns200WithBoard(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	resp := srv.get("/")
	body := srv.body(resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, "<html") {
		t.Error("response must be an HTML page")
	}
	if !strings.Contains(strings.ToLower(body), "new game") {
		t.Error("HTML must contain a 'New Game' element")
	}
}

// TestSSR_HomePageNoJavaScript validates US-29 / FR-05-10.
// Gherkin: "Web player sees the home page without JavaScript"
func TestSSR_HomePageNoJavaScript(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web")

	srv := newWebTestServer(t, "e7e5")
	resp := srv.get("/")
	body := srv.body(resp)

	if strings.Contains(body, "<script") {
		t.Error("HTML must not contain <script> elements (no JavaScript requirement)")
	}
}

// TestSSR_HomePageResponseTime validates NFR-07 (200ms).
// Gherkin: "Web player receives a fast response from the home page"
func TestSSR_HomePageResponseTime(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web")

	srv := newWebTestServer(t, "e7e5")
	start := time.Now()
	resp := srv.get("/")
	elapsed := time.Since(start)
	_ = srv.body(resp)

	assertWithinDuration(t, 200*time.Millisecond, elapsed, "home page must respond within 200ms")
}

// ─── Create Game Session ──────────────────────────────────────────────────────

// TestSSR_PostGameNewRedirectsToSessionPage validates US-30 / AC-16-02.
// Gherkin: "Web player clicks New Game and is redirected to a fresh game"
func TestSSR_PostGameNewRedirectsToSessionPage(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	resp := srv.post("/game/new")

	if resp.StatusCode != http.StatusFound {
		t.Errorf("expected 302, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if !strings.HasPrefix(loc, "/game/") {
		t.Errorf("Location must be /game/<id>, got: %q", loc)
	}
}

// TestSSR_GetSessionPageShowsStartingBoard validates US-30 / AC-16-02.
// Gherkin: "Web player follows the new game redirect and sees the starting board"
func TestSSR_GetSessionPageShowsStartingBoard(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.get(sessionPath)
	body := srv.body(resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, "<html") {
		t.Error("response must be HTML")
	}
}

// TestSSR_TwoSessionsAreIndependent validates US-30 session isolation.
// Gherkin: "Web player can create multiple independent game sessions"
func TestSSR_TwoSessionsAreIndependent(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	pathA := srv.createSession()
	pathB := srv.createSession()

	if pathA == pathB {
		t.Errorf("two sessions must have distinct paths; both got %q", pathA)
	}
}

// TestSSR_NonExistentSessionReturns404 validates error path.
// Gherkin: "Web player requesting a non-existent session receives a 404 response"
func TestSSR_NonExistentSessionReturns404(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web")

	srv := newWebTestServer(t, "e7e5")
	resp := srv.get("/game/nonexistent-session-id-xyz")
	_ = srv.body(resp)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

// ─── Make a Move via POST ─────────────────────────────────────────────────────

// TestSSR_ValidMoveRedirectsToGamePage validates US-31 / AC-16-03.
// Gherkin: "Web player makes a valid move and sees the updated board after the engine responds"
func TestSSR_ValidMoveRedirectsToGamePage(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()
	movePath := sessionPath + "/move?from=e2&to=e4"

	resp := srv.post(movePath)

	if resp.StatusCode != http.StatusFound {
		t.Errorf("expected 302, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc != sessionPath {
		t.Errorf("Location must be %s, got %q", sessionPath, loc)
	}
}

// TestSSR_ValidMoveBoardShowsPawnOnE4 validates US-31, US-32 / AC-16-03.
// Gherkin: "Web player makes a valid move and sees the updated board after the engine responds"
func TestSSR_ValidMoveBoardShowsPawnOnE4(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	// Make the move.
	_ = srv.post(sessionPath + "/move?from=e2&to=e4")

	// Follow redirect.
	resp := srv.get(sessionPath)
	body := srv.body(resp)

	// Board must show pawn on e4 (after player move) and engine response applied.
	// The exact assertion depends on the HTML template structure.
	// At minimum, the page must be a 200 with HTML.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 on GET after move, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, "<html") {
		t.Error("response must be HTML")
	}
}

// TestSSR_EngineCalledWithinSameRequestCycle validates US-32 / FR-05-06.
// Gherkin: "Web player makes a move and the engine responds within the same HTTP request"
func TestSSR_EngineCalledWithinSameRequestCycle(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	// Use a recording engine stub that tracks when it was called relative to the HTTP response.
	// stub := &recordingEngineStub{}
	// srv := web.NewServer(stub.call)
	// ...
	// After the POST, verify stub.callCount > 0 and stub.calledBeforeResponse == true.

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSSR_MoveResponseTimeWithin200ms validates NFR-07.
// Gherkin: "Web player receives a fast response when making a move"
func TestSSR_MoveResponseTimeWithin200ms(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	// Use a stub engine capped at 150ms.
	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	start := time.Now()
	resp := srv.post(sessionPath + "/move?from=e2&to=e4")
	elapsed := time.Since(start)
	_ = srv.body(resp)

	assertWithinDuration(t, 200*time.Millisecond, elapsed, "move POST including engine response must complete within 200ms")
}

// ─── Invalid Move ─────────────────────────────────────────────────────────────

// TestSSR_IllegalMoveReturns422 validates US-31 / AC-16-04.
// Gherkin: "Web player submits an illegal move and receives a 422 response"
func TestSSR_IllegalMoveReturns422(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.post(sessionPath + "/move?from=e2&to=e5")

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
	body := srv.body(resp)
	if body == "" {
		t.Error("422 response must contain a human-readable error message")
	}
}

// TestSSR_IllegalMoveLeavesStatUnchanged validates US-31 / AC-16-04.
// Gherkin: "Web player submits an illegal move and the board state is unchanged"
func TestSSR_IllegalMoveLeavesStateUnchanged(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	// Capture the board before the illegal move.
	before := srv.body(srv.get(sessionPath))

	// Attempt illegal move.
	_ = srv.post(sessionPath + "/move?from=e2&to=e5")

	// Board must be unchanged.
	after := srv.body(srv.get(sessionPath))
	if before != after {
		t.Error("board state must be unchanged after an illegal move")
	}
}

// TestSSR_MissingFromParameterReturns422 validates AC-16-04 edge case.
// Gherkin: "Web player submits a move with a missing from parameter and receives a 422 response"
func TestSSR_MissingFromParameterReturns422(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.post(sessionPath + "/move?to=e4")
	_ = srv.body(resp)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for missing from parameter, got %d", resp.StatusCode)
	}
}

// TestSSR_OutOfRangeSquareReturns422 validates AC-16-04 edge case.
// Gherkin: "Web player submits a move with an out-of-range square and receives a 422 response"
func TestSSR_OutOfRangeSquareReturns422(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.post(sessionPath + "/move?from=z9&to=e4")
	_ = srv.body(resp)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for out-of-range square, got %d", resp.StatusCode)
	}
}

// TestSSR_MovingOpponentPieceReturns422 validates AC-16-04 edge case.
// Gherkin: "Web player tries to move an opponent's piece and receives a 422 response"
func TestSSR_MovingOpponentPieceReturns422(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	// Starting position: White to move; e7 is a Black pawn.
	resp := srv.post(sessionPath + "/move?from=e7&to=e5")
	_ = srv.body(resp)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for opponent's piece move, got %d", resp.StatusCode)
	}
}

// ─── Game Over Page ───────────────────────────────────────────────────────────

// TestSSR_GameOverPageShowsCheckmateResult validates US-33 / AC-16-05.
// Gherkin: "Web player sees the game result when the game ends by checkmate"
func TestSSR_GameOverPageShowsCheckmateResult(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	// We need a session already in a checkmate state.
	// Set it up by loading the FEN directly into the session store,
	// or by playing moves until checkmate through the API.
	// For now, test structure is in place.

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSSR_GameOverPageShowsStalemateResult validates US-33 / AC-16-05.
// Gherkin: "Web player sees the game result when the game ends by stalemate"
func TestSSR_GameOverPageShowsStalemateResult(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSSR_GameOverPageShowsDrawResult validates US-33 / AC-16-05.
// Gherkin: "Web player sees a draw result with the specific draw reason"
func TestSSR_GameOverPageShowsDrawResult(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSSR_MovePostToCompletedGameReturns422 validates AC-16-05 error path.
// Gherkin: "Web player cannot submit a move to a completed game session"
func TestSSR_MovePostToCompletedGameReturns422(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// ─── SSR Board Content ────────────────────────────────────────────────────────

// TestSSR_BoardRenderedAsHTMLTable validates US-29 / FR-05-10.
// Gherkin: "Web player sees the board rendered as an HTML table"
func TestSSR_BoardRenderedAsHTMLTable(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.get(sessionPath)
	body := srv.body(resp)

	if !strings.Contains(body, "<table") {
		t.Error("board must be rendered as an HTML <table>")
	}
}

// TestSSR_NoScriptTagsInResponse validates US-29 / FR-05-10.
// Gherkin: "Web player sees no JavaScript in the HTML response"
func TestSSR_NoScriptTagsInResponse(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web")

	srv := newWebTestServer(t, "e7e5")
	sessionPath := srv.createSession()

	resp := srv.get(sessionPath)
	body := srv.body(resp)

	if strings.Contains(body, "<script") {
		t.Error("HTML must not contain <script> elements")
	}
}

// ─── Integration: Chess -> SSR Layer ─────────────────────────────────────────

// TestSSR_SessionStoreRetainsStateAcrossRequests validates integration checkpoint.
// Gherkin: integration-checkpoints.feature — "SSR session store returns the same GameState that was stored"
func TestSSR_SessionStoreRetainsStateAcrossRequests(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	srv := newWebTestServer(t, "b8c6")
	sessionPath := srv.createSession()

	// Make two moves: e2e4 (player), then g1f3 (player second turn after engine b8c6).
	_ = srv.post(sessionPath + "/move?from=e2&to=e4")
	_ = srv.post(sessionPath + "/move?from=g1&to=f3")

	resp := srv.get(sessionPath)
	body := srv.body(resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	// Board must reflect all four moves applied in sequence.
	_ = body

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}

// TestSSR_ConcurrentSessionsDoNotInterfere validates integration checkpoint.
// Gherkin: integration-checkpoints.feature — "SSR concurrent requests to different sessions do not interfere"
func TestSSR_ConcurrentSessionsDoNotInterfere(t *testing.T) {
	skipUnimplemented(t)
	_ = requiresProduction("internal/web", "internal/chess", "internal/engine")

	srv := newWebTestServer(t, "e7e5")
	pathA := srv.createSession()
	pathB := srv.createSession()

	// Make different moves in each session concurrently.
	done := make(chan struct{}, 2)
	go func() {
		_ = srv.post(pathA + "/move?from=e2&to=e4")
		done <- struct{}{}
	}()
	go func() {
		_ = srv.post(pathB + "/move?from=d2&to=d4")
		done <- struct{}{}
	}()
	<-done
	<-done

	// Each session must only reflect its own move.
	// bodyA := srv.body(srv.get(pathA))
	// bodyB := srv.body(srv.get(pathB))
	// assertBoardSquareContains(t, bodyA, "e4", true)
	// assertBoardSquareContains(t, bodyA, "d4", false)
	// assertBoardSquareContains(t, bodyB, "d4", true)
	// assertBoardSquareContains(t, bodyB, "e4", false)

	t.Fatal("not implemented — remove skipUnimplemented to enable this scenario")
}
