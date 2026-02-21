# ADR-003: SSR GUI Approach

**Status**: Accepted
**Date**: 2026-02-21
**Deciders**: Morgan (solution architect)
**Affected components**: `internal/web` (handler.go, session.go, template.go), `cmd/chess-server`

---

## Context

The SSR GUI (FR-05) must serve an interactive chess board in a browser with no JavaScript required (FR-05-10, NFR SSR constraint). The server must:
- Serve an HTML chess board at GET / (FR-05-01)
- Create sessions via POST /game/new (FR-05-02)
- Render the board at GET /game/{id} (FR-05-03)
- Accept moves via POST /game/{id}/move (FR-05-04)
- Return HTTP 422 on illegal moves (FR-05-05)
- Have the engine respond within the same request cycle as the player move (FR-05-06)
- Render in < 200ms including engine move (NFR-07)
- Store sessions in memory only (FR-05-09)
- Use no JavaScript for core functionality (FR-05-10)

The fundamental interaction model for a no-JS web application is **form POST followed by redirect** (the PRG pattern). The architecture decision is about the framework choice, template approach, session storage, and board rendering method.

---

## Decision

Use **Go standard library exclusively**: `net/http` for routing (Go 1.22 ServeMux with path parameters), `html/template` for rendering, `sync.RWMutex` + `map` for session storage, `crypto/rand` for session IDs.

**Interaction model**: Post-Redirect-Get (PRG)
- `POST /game/{id}/move` validates move, applies player move, calls engine, stores new state, returns HTTP 302
- `GET /game/{id}` renders the current state as a full HTML page
- Browser follows redirect; no double-submit on refresh

**Board rendering**: HTML `<table>` with Unicode piece characters in `<td>` elements. CSS (inline or `<style>` block) provides alternating square colors. No SVG, no image files, no external CSS framework.

**Move input**: Two `<input type="hidden">` fields inside `<form method="POST">`. A two-click selection UI requires JavaScript. Without JS, a text input for UCI-format move (e.g. "e2e4") submitted by a single form is the no-JS fallback. The form POST endpoint accepts `from` and `to` query parameters (as specified in FR-05-04), populated by the form.

**Session storage**: `map[string]chess.Game` protected by `sync.RWMutex`. Session ID is a 16-byte hex string from `crypto/rand`. No expiry in v1 (documented limitation).

**Router**: Go 1.22 `net/http.ServeMux` with wildcard patterns: `POST /game/{id}/move`, `GET /game/{id}`.

---

## Alternatives Considered

### Alternative A: Third-Party Router (chi, gorilla/mux)

**Description**: Use a well-known HTTP router package to handle path parameters and middleware.

**Evaluation against requirements**:
- Functionality: `chi` and `gorilla/mux` are more ergonomic than ServeMux for complex routing
- Dependency constraint: engine/interface layers "may use standard library and minimal well-known dependencies" (requirements.md)
- Go 1.22 change: ServeMux gained wildcard path parameters (`{id}`) in Go 1.22, eliminating the main reason to use a third-party router for simple APIs
- Routes needed: 4 routes total (`GET /`, `POST /game/new`, `GET /game/{id}`, `POST /game/{id}/move`) — well within ServeMux capability

**Rejection rationale**: The requirement specifies Go 1.22+, which provides native path parameters in ServeMux. Four routes do not justify an external router dependency. Adding `chi` or `gorilla/mux` for four routes violates the open-source-first, minimal-dependency principle.

### Alternative B: HTMX for Partial Page Updates

**Description**: Use the HTMX library (a JS-less-feeling library that uses HTML attributes) to submit forms and swap partial HTML fragments without full-page reload.

**Evaluation against requirements**:
- No-JS requirement: HTMX is a JavaScript library; it violates FR-05-10 ("no JS required for core functionality")
- UX: partial updates would improve perceived performance by not re-rendering the full page
- Dependency: adds an external JS asset dependency (CDN or bundled)

**Rejection rationale**: HTMX is JavaScript. FR-05-10 explicitly requires "no JS required for core functionality". The PRG pattern with full-page re-render satisfies the requirement. HTMX is a v2 progressive enhancement candidate.

### Alternative C: Server-Sent Events (SSE) for Engine Move Streaming

**Description**: After player move POST, immediately redirect; engine runs asynchronously; browser receives engine move via SSE and updates board with JS.

**Evaluation against requirements**:
- No-JS requirement: SSE updates require JavaScript to listen and apply DOM changes — violates FR-05-10
- Complexity: requires async session state, SSE endpoint, and JS event handler

**Rejection rationale**: SSE requires JavaScript on the client. The requirement is no-JS. The synchronous engine call within the POST request (FR-05-06 states "engine responds immediately after player move, same request cycle") is both simpler and explicitly specified.

### Alternative D: In-Memory SQLite via `modernc.org/sqlite`

**Description**: Use a pure-Go SQLite driver to persist sessions to a file-backed database.

**Evaluation against requirements**:
- Persistence requirement: FR-05-09 states "Game session persists in memory (no external DB in v1)" — explicitly in-memory
- Complexity: adds a dependency and file system coupling for no benefit in v1
- Portability: modernc.org/sqlite is CGo-free but still a significant dependency

**Rejection rationale**: Explicitly out-of-scope for v1 per FR-05-09. A `sync.RWMutex`-protected map is correct, simple, and dependency-free.

---

## Consequences

**Positive**:
- Zero external dependencies in web package (standard library only)
- PRG pattern prevents duplicate form submission on browser Back
- `html/template` is XSS-safe by default; no injection risk in board rendering
- Full-page re-render is simple and predictable; no partial-update state management
- Go 1.22 ServeMux covers routing needs natively

**Negative**:
- Full-page re-render is slower than partial updates; 200ms SSR budget requires engine to be fast
- No-JS click-to-move requires two form submissions (select piece, then select square) or a single UCI-format text input; text input is less intuitive than drag-and-drop
- In-memory sessions are lost on server restart (documented v1 limitation)
- No session expiry; a long-running server accumulates session entries indefinitely (v2: add TTL cleanup goroutine)

**200ms Budget Analysis**:
```
HTTP parse + route:   ~1ms
Session load:         ~0.1ms
Move validation:      ~1ms
Engine search:        movetime capped at 150ms in SSR mode
State store:          ~0.1ms
Template render:      ~5ms
HTTP response:        ~1ms
Total:                ~160ms (40ms headroom)
```

The SSR handler passes a `TimeControl{MoveTime: 150ms}` to the engine, reserving 50ms for HTTP overhead. This satisfies the 200ms budget while giving the engine meaningful search time.
