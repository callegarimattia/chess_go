// Package acceptance implements the executable acceptance test suite for the chess-engine feature.
//
// These tests mirror the Gherkin scenarios defined in tests/acceptance/chess-engine/*.feature.
// The .feature files are the living documentation; these _test.go files are the executable specs.
//
// Architecture alignment:
//   - chess package tests call internal/chess public API only (driving port)
//   - engine package tests call internal/engine public API only (driving port)
//   - TUI tests drive internal/tui via injected io.Reader / io.Writer
//   - SSR tests drive internal/web via net/http/httptest (driving port)
//
// Build tags:
//   - Default build (no tag): compiles but all @skip tests are skipped via t.Skip()
//   - -tags integration: enables integration checkpoint tests
//   - -tags slow: enables perft depth-5 and other long-running tests
//
// One-at-a-time workflow:
//   Remove a t.Skip() call, run `go test ./tests/acceptance/chess-engine/steps/...`,
//   implement the production code until that test passes, commit, then move to the next.

package acceptance_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	StartingFEN             = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	FoolsMateFEN            = "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
	KiwipeteFEN             = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	AfterE2E4FEN            = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	KingsOnlyFEN            = "8/8/8/8/8/8/8/K6k w - - 0 1"
	StalemateFEN            = "k7/8/1Q6/8/8/8/8/7K b - - 0 1"
	PawnOnE7FEN             = "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
	MateIn1FEN              = "k7/8/1K6/8/8/8/8/R7 w - - 0 1"
	WhiteKingsideFEN        = "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"
	CastlingThroughCheckFEN = "rnbqk2r/pppp1ppp/5n2/4p3/1b2P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 4 4"
	WhiteQueensideFEN       = "r3kbnr/ppp1pppp/2nqb3/3p4/3P4/2NQB3/PPP1PPPP/R3KBNR w KQkq - 4 5"
)

// ─── Shared Helpers ───────────────────────────────────────────────────────────

// mustBuildBinary compiles the chess-go binary and returns its path.
// The binary is written to a temp directory and cleaned up after the test.
func mustBuildBinary(t *testing.T, pkg string) string {
	t.Helper()
	dir := t.TempDir()
	binPath := dir + "/chess-go-test"
	cmd := exec.Command("go", "build", "-o", binPath, pkg)
	cmd.Dir = projectRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary %s: %v\n%s", pkg, err, out)
	}
	return binPath
}

// projectRoot walks up from the test file location to find the go.mod root.
func projectRoot(t *testing.T) string {
	t.Helper()
	// The working directory during go test is the package directory.
	// The project root contains go.mod.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}
		parent := dir[:strings.LastIndex(dir, "/")]
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// requiresProduction is called at the top of every test that exercises production code.
// It documents which production packages must be implemented before the test can pass.
// It does NOT call t.Skip() — that is done by skipUnimplemented.
func requiresProduction(packages ...string) string {
	return "requires: " + strings.Join(packages, ", ")
}

// skipUnimplemented marks a test as skipped with a clear message indicating
// that the corresponding @skip Gherkin scenario has not yet been enabled.
// Remove this call when enabling a scenario for the next TDD iteration.
func skipUnimplemented(t *testing.T) {
	t.Helper()
	t.Skip("@skip: remove this call to enable this scenario for the next TDD iteration")
}

// assertWithinDuration asserts that the elapsed time is within the allowed duration.
func assertWithinDuration(t *testing.T, allowed time.Duration, elapsed time.Duration, msgAndArgs ...interface{}) {
	t.Helper()
	if elapsed > allowed {
		t.Errorf("operation took %v, expected <= %v. %v", elapsed, allowed, fmt.Sprint(msgAndArgs...))
	}
}

// ─── Subprocess Helper ─────────────────────────────────────────────────────

// UCISession manages a UCI subprocess for acceptance tests.
type UCISession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	t      *testing.T
}

// newUCISession starts the chess-go binary as a UCI subprocess.
func newUCISession(t *testing.T, binPath string) *UCISession {
	t.Helper()
	cmd := exec.Command(binPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("StdinPipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("StdoutPipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start subprocess: %v", err)
	}
	s := &UCISession{cmd: cmd, stdin: stdin, stdout: stdout, t: t}
	t.Cleanup(func() {
		_ = stdin.Close()
		_ = cmd.Wait()
	})
	return s
}

// send writes a UCI command line to the subprocess stdin.
func (s *UCISession) send(cmd string) {
	s.t.Helper()
	if _, err := fmt.Fprintln(s.stdin, cmd); err != nil {
		s.t.Fatalf("send %q: %v", cmd, err)
	}
}

// readUntil reads stdout lines until a line containing the target string is found,
// or until the deadline is exceeded. Returns all lines read.
func (s *UCISession) readUntil(target string, deadline time.Duration) ([]string, bool) {
	s.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	lineCh := make(chan string, 64)
	go func() {
		buf := make([]byte, 4096)
		var line strings.Builder
		for {
			n, err := s.stdout.Read(buf)
			if n > 0 {
				for _, ch := range string(buf[:n]) {
					if ch == '\n' {
						lineCh <- line.String()
						line.Reset()
					} else {
						line.WriteRune(ch)
					}
				}
			}
			if err != nil {
				return
			}
		}
	}()

	var lines []string
	for {
		select {
		case <-ctx.Done():
			return lines, false
		case line := <-lineCh:
			lines = append(lines, line)
			if strings.Contains(line, target) {
				return lines, true
			}
		}
	}
}

// ─── TestMain ─────────────────────────────────────────────────────────────────

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
