# Technology Stack — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved

All technology choices default to open source. No proprietary or paid libraries are used.

---

## Language Runtime

| Component | Technology | Version | License | Rationale |
|-----------|-----------|---------|---------|-----------|
| Language | Go | 1.22+ | BSD-3-Clause | Requirement; statically typed, excellent concurrency, single-binary output |
| Module system | Go modules | built-in | BSD-3-Clause | Standard Go dependency management; no alternatives needed |

---

## chess Package (internal/chess) — Zero External Dependencies

The chess package uses only the Go standard library. No external packages are permitted.

| Concern | Standard Library Facility | Notes |
|---------|--------------------------|-------|
| FEN parsing | `strings`, `strconv` | Split on spaces and slashes |
| Error types | `errors`, `fmt` | Typed sentinel errors: ErrIllegalMove, ErrInvalidFEN |
| PGN export | `strings`, `fmt`, `time` | StringBuilder pattern for PGN text |
| Perft tests | `testing` | Table-driven tests with known perft values |
| Zobrist hashing | `math/rand` (seeded) | 64-bit hash per piece/square for repetition tracking |

---

## engine Package (internal/engine)

| Concern | Technology | License | Rationale |
|---------|-----------|---------|-----------|
| Core search | Go standard library | BSD-3-Clause | Alpha-beta is self-contained; no external search framework |
| Concurrency/cancellation | `context` (stdlib) | BSD-3-Clause | `context.WithDeadline` for time management; no goroutine leak |
| UCI I/O | `bufio`, `os` (stdlib) | BSD-3-Clause | Line-buffered stdin/stdout; no external CLI library needed |
| Time management | `time` (stdlib) | BSD-3-Clause | `time.Now()`, `time.Since()`, `time.Until()` |

---

## tui Package (internal/tui)

| Concern | Technology | License | Rationale |
|---------|-----------|---------|-----------|
| Terminal I/O | `bufio`, `os` (stdlib) | BSD-3-Clause | Sufficient for line-buffered input and ANSI output |
| Raw mode (optional) | `golang.org/x/term` | BSD-3-Clause | Enables raw keypress if needed; well-maintained official Go extension |
| Board render | `fmt`, `strings` (stdlib) | BSD-3-Clause | ASCII/Unicode output; no TUI framework needed |

**Note on `golang.org/x/term`**: This is the only external dependency in the tui package, and only if raw terminal mode is required. If the TUI uses line-by-line input (press Enter after each move), `golang.org/x/term` is not needed. Decision deferred to implementation phase.

---

## web Package (internal/web)

| Concern | Technology | License | Rationale |
|---------|-----------|---------|-----------|
| HTTP server | `net/http` (stdlib) | BSD-3-Clause | Production-grade HTTP server; no framework overhead needed for 4 routes |
| HTML templating | `html/template` (stdlib) | BSD-3-Clause | XSS-safe templating; no JS required; SSR requirement satisfied |
| Session store | `sync.RWMutex` + `map` (stdlib) | BSD-3-Clause | In-memory only (v1 requirement); no external session library |
| UUID for session IDs | `crypto/rand` (stdlib) | BSD-3-Clause | `crypto/rand` produces session IDs without external library |
| Router | `net/http.ServeMux` (stdlib) | BSD-3-Clause | Go 1.22 enhanced ServeMux supports `{id}` path parameters natively |

**Note on Go 1.22 ServeMux**: Go 1.22 added wildcard pattern matching (`/game/{id}`) directly in `net/http.ServeMux`. This eliminates the need for a third-party router (gorilla/mux, chi). The 1.22+ language requirement is intentional and enables this.

---

## Testing

| Concern | Technology | License | Rationale |
|---------|-----------|---------|-----------|
| Unit tests | `testing` (stdlib) | BSD-3-Clause | Standard Go test runner; table-driven tests |
| Perft tests | `testing` (stdlib) | BSD-3-Clause | Known position node counts as test cases |
| HTTP handler tests | `net/http/httptest` (stdlib) | BSD-3-Clause | In-process HTTP testing without network |
| Benchmarks | `testing.B` (stdlib) | BSD-3-Clause | NPS benchmarks using Go benchmark framework |
| Race detector | `go test -race` (stdlib) | BSD-3-Clause | Concurrent GameState access validation |

---

## Build and Tooling

| Concern | Technology | License | Rationale |
|---------|-----------|---------|-----------|
| Build | `go build` | BSD-3-Clause | Standard; no build system needed for two binaries |
| Test | `go test` | BSD-3-Clause | Standard test runner |
| Lint | `golangci-lint` | MIT | Aggregates staticcheck, errcheck, vet; industry standard |
| Format | `gofmt` / `goimports` | BSD-3-Clause | Standard Go formatting |
| CI | GitHub Actions | MIT (runner) | Free for public repos; YAML-configured |

---

## Rejected Technologies

| Technology | Reason for Rejection |
|-----------|---------------------|
| `github.com/notnil/chess` | Defeats the purpose; goal is to build, not import a chess library |
| PostgreSQL / SQLite | No persistence requirement in v1; in-memory session store is sufficient |
| React / Vue / HTMX | SSR requirement explicitly excludes JavaScript; html/template is sufficient |
| Gorilla/mux | Go 1.22 ServeMux covers the routing needs natively |
| Cobra (CLI framework) | Two simple binaries need no CLI framework; adds dependency without benefit |
| Bubble Tea / Lip Gloss | Full TUI framework is overengineered; ASCII render with fmt is sufficient for v1 |
| gRPC | Internal function calls between packages; no RPC needed in monolith |

---

## License Summary

| Package | License | Source |
|---------|---------|--------|
| Go standard library | BSD-3-Clause | https://go.dev |
| golang.org/x/term | BSD-3-Clause | https://pkg.go.dev/golang.org/x/term |
| golangci-lint | MIT | https://github.com/golangci/golangci-lint |

All licenses are permissive (BSD-3-Clause, MIT). No GPL, AGPL, or proprietary licenses.
