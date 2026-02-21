# Observability Design — chess-engine
**Feature**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved
**Author**: Apex (nw-platform-architect)

---

## 1. Observability Philosophy

This is a local CLI tool and a local HTTP server. There is no production infrastructure to instrument, no external logging sink, and no dashboards. The observability model is:

- **Structured stderr logs** in both binaries for developer diagnostics
- **Go benchmark suite** in the repository for performance validation
- **No external infrastructure** required (no Prometheus, no ELK, no Datadog)

Observability serves two audiences:
1. **Daniel (engine developer)**: Debug search behaviour, validate NPS, diagnose time management
2. **CI pipeline**: Benchmark regression detection, perft correctness validation

---

## 2. Structured Logging

### 2.1 Log Format

All log output uses structured JSON to stderr. This enables `jq` filtering in development without requiring a logging service.

```json
{"time":"2026-02-21T14:32:01.234Z","level":"info","component":"engine","msg":"search complete","depth":6,"score_cp":32,"nodes":143820,"nps":119850,"duration_ms":1200,"best_move":"g1f3"}
```

### 2.2 Log Fields (Mandatory)

| Field | Type | Description | Example |
|---|---|---|---|
| `time` | string (RFC3339) | Timestamp with millisecond precision | `"2026-02-21T14:32:01.234Z"` |
| `level` | string | Log level: `debug`, `info`, `warn`, `error` | `"info"` |
| `component` | string | Package emitting the log | `"engine"`, `"tui"`, `"web"` |
| `msg` | string | Human-readable description | `"search complete"` |

### 2.3 Component-Specific Fields

**engine component**:

| Field | Type | Description |
|---|---|---|
| `depth` | int | Search depth completed |
| `score_cp` | int | Centipawn score (positive = white advantage) |
| `nodes` | int | Total nodes evaluated |
| `nps` | int | Nodes per second |
| `duration_ms` | int | Search duration in milliseconds |
| `best_move` | string | UCI move string (e.g., `"g1f3"`) |

Example log sequence during search:
```json
{"time":"...","level":"debug","component":"engine","msg":"search start","fen":"rnbqkbnr/.../w KQkq - 0 1","movetime_ms":1000}
{"time":"...","level":"debug","component":"engine","msg":"depth complete","depth":1,"score_cp":0,"nodes":20,"nps":200000,"duration_ms":0}
{"time":"...","level":"debug","component":"engine","msg":"depth complete","depth":2,"score_cp":15,"nodes":400,"nps":180000,"duration_ms":2}
{"time":"...","level":"info","component":"engine","msg":"search complete","depth":6,"score_cp":32,"nodes":143820,"nps":119850,"duration_ms":1200,"best_move":"g1f3"}
```

**tui component**:

| Field | Type | Description |
|---|---|---|
| `move` | string | Move applied (UCI format) |
| `result` | string | Game result if terminal (`"checkmate"`, `"stalemate"`, `"draw"`) |

Example:
```json
{"time":"...","level":"info","component":"tui","msg":"player move applied","move":"e2e4"}
{"time":"...","level":"info","component":"tui","msg":"engine move applied","move":"e7e5","duration_ms":450}
{"time":"...","level":"info","component":"tui","msg":"game ended","result":"checkmate","winner":"white","moves":40}
```

**web component**:

| Field | Type | Description |
|---|---|---|
| `method` | string | HTTP method |
| `path` | string | Request path |
| `status` | int | HTTP response status code |
| `duration_ms` | int | Request duration in milliseconds |
| `session_id` | string | First 8 chars of session ID (truncated for readability) |

Example:
```json
{"time":"...","level":"info","component":"web","msg":"request","method":"POST","path":"/game/abc123/move","status":302,"duration_ms":148,"session_id":"abc12345"}
{"time":"...","level":"warn","component":"web","msg":"illegal move","method":"POST","path":"/game/abc123/move","status":422,"duration_ms":2,"session_id":"abc12345","move_from":"e2","move_to":"e9"}
```

### 2.4 Log Levels

| Level | Use |
|---|---|
| `debug` | Per-depth search info; individual node counts; verbose engine state |
| `info` | Search completion; player/engine moves; game results; HTTP requests |
| `warn` | Illegal move attempts; UCI unknown commands; unexpected but handled conditions |
| `error` | FEN parse failure; I/O errors; any condition requiring user attention |

Default level: `info`. Controlled by `LOG_LEVEL` environment variable or `-log-level` flag (future).

### 2.5 Output Routing

| Binary | User-visible output | Log output |
|---|---|---|
| `chess-go` (TUI) | Board, status messages, prompts → `stdout` | Engine/TUI diagnostics → `stderr` |
| `chess-server` (SSR) | HTTP responses (no stdout) | Request logs, engine diagnostics → `stderr` |

This separation ensures TUI board output is clean (no JSON noise for Sofia) while engine logs remain capturable by Daniel for debugging.

UCI protocol output (for `cmd/chess-go` in UCI mode) goes to `stdout` as required by the UCI spec. UCI info lines double as observability for Daniel in GUI mode.

### 2.6 Implementation Notes

Use Go's `log/slog` package (stdlib, available since Go 1.21; project uses 1.22+). `slog.JSONHandler` writing to `os.Stderr` satisfies all requirements with zero external dependencies.

```go
// Recommended initialisation in main()
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelInfo, // configurable via LOG_LEVEL
}))
slog.SetDefault(logger)
```

---

## 3. Performance Observability — Benchmark Suite

### 3.1 Benchmark File: `internal/engine/bench_test.go`

The benchmark suite validates NFR-01 (>100,000 NPS) and detects regressions in CI.

**Required benchmarks**:

```
BenchmarkSearchDepth4    — Search from starting position to depth 4
BenchmarkSearchDepth5    — Search from starting position to depth 5
BenchmarkSearchDepth6    — Search from starting position to depth 6; depth 6 is the ceiling for current representation
BenchmarkMoveGeneration  — LegalMoves() from starting position
BenchmarkMoveApplication — Apply() for 20 moves from starting position
BenchmarkEvaluation      — Eval() at a mid-game position
```

**NPS extraction**: The benchmark records `b.N` iterations and elapsed time. NPS is derived as:
```
NPS = nodes_searched / elapsed_seconds
```

Each `BenchmarkSearchDepthN` counts nodes via `SearchResult.Nodes` and reports via `b.ReportMetric(float64(nodes)/elapsed, "nps")`.

### 3.2 Benchmark Positions

| Position | FEN | Purpose |
|---|---|---|
| Starting | `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1` | Baseline NPS measurement |
| Kiwipete | `r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1` | Complex mid-game; higher branching factor |
| End-game | `8/8/8/8/8/8/4K3/4k3 w - - 0 1` | Low branching factor; deep search test |

### 3.3 Regression Detection in CI

CI Stage 7 procedure:
1. Run `go test -bench=. -benchmem -count=3 -benchtime=3s ./... > bench-current.txt`
2. If `bench-baseline.txt` exists in cache: `benchstat bench-baseline.txt bench-current.txt`
3. Parse benchstat output; fail CI if any benchmark shows `+delta > 20%` with p < 0.05
4. On main branch success: cache `bench-current.txt` as new `bench-baseline.txt`

`benchstat` is the standard Go benchmark analysis tool from `golang.org/x/perf`. It applies statistical significance testing to avoid flaky regressions on noisy CI runners.

### 3.4 Perft Correctness Suite

Located in `internal/chess/perft_test.go` (build tag `slow`):

**Known perft values (starting position)**:

| Depth | Nodes |
|---|---|
| 1 | 20 |
| 2 | 400 |
| 3 | 8,902 |
| 4 | 197,281 |
| 5 | 4,865,609 |

**Kiwipete perft values**:

| Depth | Nodes |
|---|---|
| 1 | 48 |
| 2 | 2,039 |
| 3 | 97,862 |
| 4 | 4,085,603 |

CI runs perft at depths 1-4 on the starting position and depths 1-3 on Kiwipete. Depth 5 requires the `-tags slow` build tag and is excluded from standard CI to keep the pipeline under 15 minutes.

---

## 4. What Observability Does Not Include

The following were explicitly excluded due to the local tool deployment model:

| Excluded | Reason |
|---|---|
| Prometheus metrics endpoint | No long-running server in production context |
| Distributed tracing (OpenTelemetry) | Single binary, single process; no service graph |
| Log aggregation (ELK, Loki) | Local tool; developer reads stderr directly |
| Error tracking (Sentry) | Local tool; no user telemetry |
| Uptime monitoring | No always-on server |
| SLO burn rate alerts | No production SLOs for a local CLI tool |

These omissions are appropriate and intentional. If `chess-server` is ever deployed as a shared service (v2+), this observability design provides the log structure foundation for migration to an aggregation service without code changes (JSON logs are already compatible with Loki, CloudWatch, and Datadog).
