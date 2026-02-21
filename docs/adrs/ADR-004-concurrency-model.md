# ADR-004: Concurrency Model

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Morgan (solution architect)
**Affected components**: `internal/engine` (search.go, uci.go), `internal/web` (session.go), `internal/tui` (loop.go)

---

## Context

The system has three concurrent execution contexts:

1. **Engine search**: Alpha-beta search is CPU-intensive and must be interruptible. The UCI `stop` command must cause the engine to emit `bestmove` within 100ms (FR-03-08). The engine must never block past `movetime + 50ms` (NFR-05).
2. **UCI handler**: Reads stdin line-by-line while the engine may be searching in a goroutine. Commands arriving during search (`stop`, `quit`) must be processed immediately.
3. **SSR session store**: Multiple concurrent HTTP requests (one per browser tab or load test) may read and write the session map simultaneously.

Key constraints:
- `GameState` is immutable (NFR-05): "GameState is immutable and goroutine-safe"
- Engine goroutine must not leak when `stop` or `quit` is received
- SSR web handler is called by `net/http`'s built-in goroutine-per-request model
- TUI game loop is single-user; minimal concurrency required

---

## Decision

### Decision 1: Engine Cancellation via `context.Context`

The search is launched in a goroutine with a `context.WithDeadline` context. The search goroutine checks `ctx.Done()` at every node. On cancellation, it returns the best move found in the last completed iteration to a result channel.

```
Caller (UCI handler or TUI loop):
  ctx, cancel = context.WithDeadline(parent, deadline)
  resultCh = make(chan SearchResult, 1)
  go func() { resultCh <- search(ctx, position, alpha, beta) }()
  select {
    case result = <-resultCh: // search completed within deadline
    case <-stopCh:            // stop command received
      cancel()
      result = <-resultCh    // wait for goroutine to exit cleanly
  }
```

No goroutines are leaked: the search goroutine always exits by sending to resultCh after ctx.Done() returns, and resultCh is buffered (capacity 1) so the send never blocks.

### Decision 2: UCI Handler — Two-Goroutine Architecture

The UCI binary uses two goroutines:

- **Reader goroutine**: Reads stdin line by line, sends parsed commands to a command channel (`chan UCICommand`). Runs independently of search.
- **Dispatcher (main goroutine)**: Selects on command channel. On `go` command, launches search goroutine with context. On `stop`, cancels context and waits for bestmove.

```
Reader goroutine: stdin → bufio.Scanner → commandCh
Main goroutine:
  select {
    case cmd := <-commandCh:
      switch cmd.Type {
        case UCI, IsReady, UCINewGame: respond immediately
        case Position: update position state
        case Go: launch search goroutine, record cancel func
        case Stop: invoke cancel func, await bestmove from resultCh
        case Quit: invoke cancel func, os.Exit(0)
      }
  }
```

Output serialization: all writes to stdout happen in the main goroutine (or via a dedicated writer goroutine with a channel). No concurrent writes to stdout — UCI requires ordered output.

### Decision 3: SSR Session Store — sync.RWMutex

The session store wraps a `map[string]Session` with a `sync.RWMutex`. Read operations (GET requests) use `RLock`; write operations (POST move, POST new game) use `Lock`.

```
type SessionStore struct {
    mu       sync.RWMutex
    sessions map[string]Session
}
```

Since `chess.Game` (SA-04) is immutable, storing and reading it from the map is race-free: writers replace the value atomically under Lock; readers get a consistent copy under RLock.

### Decision 4: TUI Game Loop — Single Goroutine with Engine Goroutine

The TUI main goroutine handles all I/O and rendering. The engine is called in a separate goroutine to allow the "Engine thinking..." status to be displayed while the search runs.

```
Main goroutine:
  render(gameState)
  print("Your move: ")
  move = readInput()
  gameState = chess.Apply(move)
  render(gameState)
  print("Engine thinking...")
  resultCh = make(chan chess.Move, 1)
  go func() { resultCh <- engine.Search(gameState, tc).BestMove }()
  engineMove = <-resultCh
  gameState = chess.Apply(engineMove)
  // loop
```

No goroutine leak: the search goroutine always sends to the buffered resultCh. The main goroutine always reads from resultCh before looping.

---

## Alternatives Considered

### Alternative A: Parallel Search (Multiple Search Goroutines)

**Description**: Launch N goroutines searching different sub-trees in parallel (Lazy SMP or YBWC). Standard in high-performance engines for multi-core utilization.

**Evaluation against requirements**:
- Performance: 100k NPS target is achievable on a single core with the 8x8 array representation
- Complexity: Lazy SMP requires a shared transposition table with atomic operations; significant implementation complexity
- Correctness: parallel search introduces non-determinism that complicates testing and debugging

**Rejection rationale**: Single-core search exceeds NFR-01 (100k NPS). Parallel search is a v2 optimization. The 8x8 board representation is the current bottleneck, not parallelism. Parallel search without a transposition table provides minimal benefit.

### Alternative B: Single Goroutine with Polling for Stop

**Description**: Run search in the main goroutine; poll for stop by checking a shared flag (atomic bool) on every node.

**Evaluation against requirements**:
- Correctness: reading stdin requires a separate goroutine; polling a flag requires the search to yield to the scheduler
- Go concurrency: Go goroutines are cooperatively scheduled; a tight CPU loop may not yield for stdin reads without a goroutine
- Complexity: atomic flag is simpler than context; but the stdin reader still needs its own goroutine

**Rejection rationale**: `context.Context` is the idiomatic Go cancellation mechanism. It integrates with the standard library, is composable with deadlines, and is well-understood by Go developers. An atomic flag provides no advantage and breaks the idiom.

### Alternative C: sync/atomic for Session Store (Lock-Free)

**Description**: Use `sync/atomic.Value` to store and swap the entire session map as a single atomic pointer swap.

**Evaluation against requirements**:
- Correctness: atomic swap of an entire map is correct for single-writer scenarios; with multiple concurrent writes, it introduces lost-update races
- Complexity: lock-free programming requires careful reasoning about memory ordering
- Actual concurrency level: SSR is local-only in v1 (FR-05-09 note); concurrent request rate is expected to be low (single human player)

**Rejection rationale**: `sync.RWMutex` is the standard Go primitive for concurrent map access. Lock-free approaches are appropriate when mutex contention is measured to be a bottleneck. At the expected concurrency level (single local user), RWMutex adds negligible overhead.

---

## Consequences

**Positive**:
- `context.Context` cancellation is idiomatic Go; composable with timeouts and parent cancellation
- Two-goroutine UCI model is simple and straightforward: one reader, one dispatcher
- `sync.RWMutex` allows concurrent reads (multiple GET requests) without blocking
- Immutable `chess.Game` eliminates all data races on game state; no locks needed on the state itself
- All goroutines exit cleanly; no goroutine leaks in normal or error paths

**Negative**:
- Context check adds a branch at every node (~1ns overhead); acceptable at 100k NPS
- UCI writer serialization in main goroutine means `info` lines from search must be sent via channel (or search writes directly to stdout with synchronization); crafter must choose approach
- RWMutex adds locking overhead for session reads; acceptable for local-only usage

**Race Condition Validation**:
- Run `go test -race ./...` as mandatory CI check
- The race detector will catch any missed synchronization on session store or engine result channel
- Immutable `chess.Game` means the detector will not flag reads of game state across goroutines

**Goroutine Leak Prevention Checklist**:
- Search goroutine: always sends to buffered result channel before returning; caller always reads from channel
- UCI reader goroutine: exits when stdin is closed (EOF) or `quit` is processed (os.Exit)
- TUI engine goroutine: always sends to buffered result channel; main loop always reads
- SSR request goroutines: managed by `net/http` server; engine call is synchronous within handler; no goroutine launched
