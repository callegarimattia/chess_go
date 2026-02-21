# Component Boundaries — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved

---

## Boundary Definition Principles

1. **Dependency direction is strict**: outer packages import inner packages; inner packages never import outer packages.
2. **chess package is the domain core**: it has no knowledge of engine, TUI, or web.
3. **engine package is the application layer**: it depends only on chess; it has no knowledge of TUI or web rendering.
4. **tui and web are adapters**: they translate between user input formats and the chess/engine API.
5. **cmd binaries are entry points**: they wire adapters to packages but contain no business logic.

---

## Component: chess package (`internal/chess`)

### Responsibility
Complete chess rules implementation as a reusable Go library. Pure domain model.

### Owns
- FEN string parsing and validation
- Board representation (piece placement, square indexing)
- GameState as an immutable value (ActiveColor, CastlingRights, EnPassantSquare, HalfMoveClock, FullMoveNumber)
- Legal move generation (all piece types, including special moves)
- Move application: returns a new GameState, never mutates
- Check, checkmate, stalemate detection
- Draw detection: fifty-move rule, threefold repetition (via position hash history), insufficient material
- Move notation: UCIString(), SANString() — output formatting only
- PGN export: Game.ToPGN()
- Typed error values: ErrIllegalMove, ErrInvalidFEN, ErrInvalidMoveFormat

### Does Not Own
- Search or evaluation (engine concern)
- Rendering (tui/web concern)
- UCI protocol (engine concern)
- Time management (engine concern)

### Dependency Rule
- **Imports**: nothing outside the Go standard library
- **Imported by**: engine, tui, web

### Public Surface (Ports)
- `GameState` struct — fields: Board, ActiveColor, CastlingRights, EnPassantSquare, HalfMoveClock, FullMoveNumber
- `Move` struct — fields: From, To, Promotion (all value types, no pointers)
- `Game` struct — wraps GameState with move history for PGN and repetition tracking
- `NewGameFromFEN(fen string) (Game, error)` — primary constructor
- `Game.LegalMoves() []Move` — complete legal move list
- `Game.Apply(m Move) (Game, error)` — immutable state transition
- `Game.InCheck() bool` — active color king attacked
- `Game.Result() GameResult` — terminal state or InProgress
- `Game.ToFEN() string` — FEN serialization
- `Game.ToPGN() string` — PGN serialization
- `Move.UCIString() string` — "e2e4", "e7e8q"
- `Move.SANString(g Game) string` — "Nf3", "O-O", "e8=Q+"

### Constraint: Immutability
Every `Game.Apply()` call returns a new `Game`. The original is unchanged. This makes `GameState` safe to share across goroutines without locks. The engine exploits this during search (multiple goroutines can hold references to the same state safely).

---

## Component: engine package (`internal/engine`)

### Responsibility
Chess search, evaluation, time management, and UCI protocol I/O.

### Owns
- Alpha-beta search with iterative deepening
- Move ordering (MVV-LVA for captures, killer move heuristic)
- Quiescence search (extends search at tactical positions)
- Material evaluation (standard piece values)
- Positional evaluation (piece-square tables, per-piece)
- Time allocation: `movetime`, `wtime/btime/winc/binc` strategies
- Context-based cancellation: search exits cleanly within 50ms of deadline
- UCI stdin/stdout protocol handling (all required commands)
- UCI info line emission during search

### Does Not Own
- Chess rules (chess package concern)
- Rendering (tui/web concern)
- HTTP or terminal I/O (tui/web concern)

### Dependency Rule
- **Imports**: `internal/chess`, Go standard library
- **Imported by**: tui, web, cmd binaries

### Public Surface (Ports)
- `Search(g chess.Game, tc TimeControl, info io.Writer) SearchResult` — primary search entry point
- `SearchResult` struct — fields: BestMove (chess.Move), Score (int centipawns), Depth (int), Nodes (int)
- `TimeControl` struct — fields: MoveTime, WTime, BTime, WInc, BInc (all time.Duration)
- `UCIHandler` struct — `Run(r io.Reader, w io.Writer)` reads commands and writes responses
- `NewUCIHandler(searchFn SearchFunc) UCIHandler` — constructor with search dependency injection

### Constraint: Time Compliance
Search MUST return within `TimeControl.MoveTime + 50ms`. The time manager sets a `context.WithDeadline` and the search goroutine respects `ctx.Done()` at the top of each node. If no move has been searched (pathological case), the first legal move is returned immediately.

---

## Component: tui package (`internal/tui`)

### Responsibility
Terminal game loop, ASCII board rendering, move input parsing, status display.

### Owns
- ASCII board rendering from GameState (8x8 grid with coordinates)
- Unicode piece symbols or ASCII fallback (configurable at compile time)
- Move prompt and input reading (line-buffered from io.Reader)
- Status line: "White/Black to move", "Engine thinking...", "Check!", result messages
- Game loop: player turn → engine turn → result check → loop or exit
- PGN save prompt at game end

### Does Not Own
- Chess rules (chess package concern)
- Engine search (engine package concern)
- HTTP (web concern)

### Dependency Rule
- **Imports**: `internal/chess`, `internal/engine`, Go standard library, optional `golang.org/x/term`
- **Imported by**: `cmd/chess-go`

### Public Surface
- `NewGame(r io.Reader, w io.Writer, engineFn EngineFunc) Game` — configurable I/O for testing
- `Game.Run()` — starts the interactive game loop
- `EngineFunc` type alias: `func(g chess.Game, tc engine.TimeControl) chess.Move`

### Design Notes
- I/O is injected (io.Reader, io.Writer) to enable testing without a real terminal
- The engine is called as a function (not a struct), enabling substitution of random-move skeleton in tests

---

## Component: web package (`internal/web`)

### Responsibility
HTTP handlers, HTML template rendering, in-memory session management for SSR GUI.

### Owns
- HTTP route handlers for: `GET /`, `POST /game/new`, `GET /game/{id}`, `POST /game/{id}/move`
- HTML templates (board as HTML table, status messages, New Game button)
- In-memory session store: map[sessionID]chess.Game protected by sync.RWMutex
- Session ID generation (crypto/rand hex string)
- Server-side move validation (delegate to chess package)
- Engine call after player move (synchronous, within HTTP request)
- HTTP 422 response on illegal moves with human-readable error body

### Does Not Own
- Chess rules (chess package concern)
- Engine search (engine package concern)
- Terminal I/O (tui concern)

### Dependency Rule
- **Imports**: `internal/chess`, `internal/engine`, Go standard library
- **Imported by**: `cmd/chess-server`

### Public Surface
- `NewServer(engineFn EngineFunc) *Server` — returns configured http.Handler
- `Server.ServeHTTP(w http.ResponseWriter, r *http.Request)` — satisfies http.Handler

### Design Notes
- No JavaScript in any template; all interactivity via form POST and redirect (PRG pattern)
- Board rendered as HTML `<table>` with piece Unicode characters in `<td>` cells
- Move input: two hidden form fields (from, to) submitted via form POST; no JS drag-and-drop in v1
- The PRG (Post-Redirect-Get) pattern prevents duplicate form submission on browser refresh

---

## Component: cmd/chess-go (`cmd/chess-go/main.go`)

### Responsibility
Entry point for TUI binary. Wire tui, engine, and I/O.

### Owns
- `os.Stdin` / `os.Stdout` binding to tui.NewGame()
- Default TimeControl for engine (configurable via flags in future)
- Process exit code management

### Does Not Own
- Game logic, rendering, or search (all delegated)

---

## Component: cmd/chess-server (`cmd/chess-server/main.go`)

### Responsibility
Entry point for SSR HTTP server. Wire web package and start HTTP listener.

### Owns
- Port binding (default :8080, configurable via env or flag)
- HTTP server lifecycle (graceful shutdown on SIGINT)
- web.NewServer() instantiation

### Does Not Own
- Request handling, chess logic, or rendering (all delegated)

---

## Inter-Component Data Contracts

| Data | Type | Producer | Consumers | Crossing Boundary |
|------|------|----------|-----------|-------------------|
| FEN string (SA-01) | string | chess.Game.ToFEN() | engine, tui, web | chess → engine, tui, web |
| Move (UCI) (SA-02) | chess.Move | tui input, web POST, engine.Search | chess.Game.Apply | tui/web → chess |
| GameState (SA-04) | chess.Game | chess.Game.Apply | engine.Search, tui.Renderer, web.Handler | chess → engine, tui, web |
| SearchResult (SA-06) | engine.SearchResult | engine.Search | tui.loop, web.handler | engine → tui, web |
| TimeControl (SA-07) | engine.TimeControl | tui.loop, web.handler | engine.Search | tui/web → engine |
| SAN notation (SA-03) | string | chess.Move.SANString | tui.Renderer, web.template | chess → tui, web |
| PGN record (SA-05) | string | chess.Game.ToPGN | filesystem, tui export | chess → tui |

---

## Boundary Enforcement (Go Mechanism)

Go's package system enforces these boundaries at compile time:
- `internal/` directory prevents any external Go module from importing internal packages
- Within the module, import direction is enforced by convention and validated by `go vet` + dependency analysis tools
- Circular imports are rejected by the Go compiler
- No global variables cross package boundaries; all state flows through function parameters and return values

---

## Anti-Patterns Explicitly Prohibited

| Anti-Pattern | Why Prohibited |
|-------------|----------------|
| chess package importing engine | Violates dependency direction; chess is innermost layer |
| Global mutable game state | Violates immutability requirement; breaks goroutine safety |
| Shared `GameState` pointer mutation | Engine search may hold references across goroutines; mutation causes data races |
| Direct FEN string manipulation in engine | chess package owns FEN; engine uses chess.Game API |
| HTML generation in chess or engine packages | Rendering is an adapter concern (web package) |
| UCI parsing in tui or web | UCI is an engine adapter concern |
