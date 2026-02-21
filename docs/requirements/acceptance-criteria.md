# Acceptance Criteria — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21

All criteria are written in Given-When-Then (GWT) format and are directly testable.
Each criterion maps to a user story and a Gherkin scenario in `journey-chess-engine.feature`.

---

## AC-00: Walking Skeleton

**Story**: US-00

**AC-00-01**: Given the starting position, when the skeleton runs end-to-end, then:
- A legal move is selected (verifiable against the 20 legal moves from the starting position)
- The updated GameState is displayed as an ASCII board in the terminal
- The process exits with code 0
- No panic or runtime error occurs

---

## AC-01: FEN Parser

**Story**: US-01

**AC-01-01**: Given the string `"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"`, when `NewGameFromFEN` is called, then:
- Active color is White
- Castling rights include all four (KQkq)
- En passant square is none
- Half-move clock is 0
- Full-move number is 1
- All 32 pieces are in starting positions

**AC-01-02**: Given any malformed FEN string (wrong field count, invalid piece chars, invalid castling string), when `NewGameFromFEN` is called, then:
- An `ErrInvalidFEN` typed error is returned
- No GameState is returned
- No panic occurs

**AC-01-03**: Given a valid FEN, when `GameState.ToFEN()` is called, then the output string round-trips: `NewGameFromFEN(gs.ToFEN())` produces an equal GameState.

---

## AC-02: Legal Move Generation

**Story**: US-02

**AC-02-01**: Given the starting position, when `LegalMoves()` is called, then exactly 20 moves are returned (16 pawn moves + 4 knight moves).

**AC-02-02**: Given any position, when `LegalMoves()` is called, then no illegal move is included (verified by attempting each move and confirming the resulting position is not self-in-check).

**AC-02-03**: Given a position in checkmate, when `LegalMoves()` is called, then an empty slice is returned.

**AC-02-04**: Given a position in stalemate, when `LegalMoves()` is called, then an empty slice is returned.

---

## AC-03: Apply Move

**Story**: US-03

**AC-03-01**: Given the starting position and the move `e2e4`, when `Apply(move)` is called, then:
- A new GameState is returned (original is unchanged)
- Active color switches to Black
- En passant square is e3
- The pawn is on e4
- The pawn is no longer on e2

**AC-03-02**: Given a position and an illegal move (not in `LegalMoves()`), when `Apply(move)` is called, then `ErrIllegalMove` is returned and the original GameState is unchanged.

---

## AC-04: Check and Checkmate Detection

**Story**: US-04

**AC-04-01**: Given the Fool's Mate position `"rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"`, when `Result()` is called, then:
- Result is `Checkmate`
- The losing side is White

**AC-04-02**: Given a position where the active color's king is attacked, when `InCheck()` is called, then `true` is returned.

**AC-04-03**: Given a position where the active color's king is not attacked, when `InCheck()` is called, then `false` is returned.

---

## AC-05: Draw Detection

**Story**: US-05

**AC-05-01**: Given a position where the half-move clock is 100, when `Result()` is called, then the result is `DrawFiftyMove`.

**AC-05-02**: Given a game where the same position has occurred three times (by repetition), when `Result()` is called, then the result is `DrawThreefoldRepetition`.

**AC-05-03**: Given a position with only kings remaining, when `Result()` is called, then the result is `DrawInsufficientMaterial`.

**AC-05-04**: Given a stalemate position (no legal moves, not in check), when `Result()` is called, then the result is `Stalemate`.

---

## AC-06: Castling

**Story**: US-06

**AC-06-01**: Given a position where White can castle kingside, when `Apply("e1g1")` is called, then:
- White king is on g1
- White rook is on f1
- Castling rights no longer include White kingside

**AC-06-02**: Given a position where the castling path passes through an attacked square, when `Apply("e1g1")` is called, then `ErrIllegalMove` is returned.

**AC-06-03**: Given a position where the king is in check, when any castling move is attempted, then `ErrIllegalMove` is returned.

**AC-06-04**: Castling moves are generated in `LegalMoves()` if and only if:
- King and rook have not moved
- Squares between king and rook are empty
- King does not pass through or land on an attacked square
- King is not currently in check

---

## AC-07: En Passant

**Story**: US-07

**AC-07-01**: Given a position where en passant is available on e6 (White pawn on d5, Black pawn just moved to e5), when `Apply("d5e6")` is called, then:
- White pawn is on e6
- Black pawn on e5 is removed
- En passant square in the new state is none

**AC-07-02**: Given a position with an en passant pin (capturing en passant would expose the king), the en passant capture is absent from `LegalMoves()`.

---

## AC-08: Pawn Promotion

**Story**: US-08

**AC-08-01**: Given a position with a White pawn on e7, when `Apply("e7e8q")` is called, then:
- A White queen is on e8
- No pawn is on e7 or e8

**AC-08-02**: Given a pawn on the 7th rank, `LegalMoves()` returns four promotion moves for each available target square (q, r, b, n).

**AC-08-03**: Applying `"e7e8"` (no promotion piece) when a pawn is on e7 returns `ErrIllegalMove`.

---

## AC-09: Notation Export

**Story**: US-09, US-10

**AC-09-01**: `Move.UCIString()` returns the correct string for all move types including promotion (`e7e8q`) and castling (`e1g1`).

**AC-09-02**: `Move.SANString(gs)` returns correct SAN for:
- Pawn moves: `e4`
- Piece moves: `Nf3`
- Captures: `Bxe5`
- Check: `Nf3+`
- Checkmate: `Qh5#`
- Castling: `O-O`, `O-O-O`
- Promotion: `e8=Q`
- Disambiguation: `Rfe1` (when two rooks can move to e1)

---

## AC-10: PGN Export

**Story**: US-11

**AC-10-01**: Given a complete game, `Game.ToPGN()` produces a string containing:
- Seven tag pairs: Event, Site, Date, Round, White, Black, Result
- Move text in SAN format with move numbers
- The result token at the end (1-0, 0-1, or 1/2-1/2)

**AC-10-02**: The produced PGN is parseable by standard PGN readers (structure validated against PGN standard).

---

## AC-11: Perft Validation

**Story**: US-12

**AC-11-01**: Starting from the initial position, perft results match known values:

| Depth | Expected nodes |
|-------|----------------|
| 1     | 20             |
| 2     | 400            |
| 3     | 8,902          |
| 4     | 197,281        |
| 5     | 4,865,609      |

**AC-11-02**: Perft from position 2 (Kiwipete) matches known values at depth 1–4.

---

## AC-12: Engine Search

**Story**: US-14

**AC-12-01**: Given any legal position and a movetime of 1000ms, when `SearchBestMove` is called, then:
- A legal bestmove is returned within 1050ms
- At least one `info depth` line was emitted

**AC-12-02**: Given the starting position with movetime 100ms, the engine searches at least depth 3.

**AC-12-03**: Given a position with a forced mate in 1, the engine finds the mating move at any search depth >= 1.

---

## AC-13: Time Management

**Story**: US-17

**AC-13-01**: Given `go movetime 500`, the engine emits `bestmove` within 550ms (500 + 50ms grace).

**AC-13-02**: Given `go wtime 60000 btime 60000` (1 minute per side), the engine does not use more than `60000 / 30` ms on average per move over a 30-move game.

---

## AC-14: UCI Protocol

**Story**: US-20–US-23

**AC-14-01**: Given the engine binary is started and `uci` is sent, then:
- `id name chess-go` is output
- `id author <non-empty>` is output
- `uciok` is output
- Response arrives within 100ms

**AC-14-02**: Given `isready` is sent after `uci`, then `readyok` is output within 100ms.

**AC-14-03**: Given `position startpos moves e2e4 e7e5` then `go movetime 500`, then `bestmove` is emitted and the move is legal in the position after the two moves.

**AC-14-04**: Given the engine is searching and `stop` is sent, then `bestmove` is emitted within 100ms of receiving `stop`.

**AC-14-05**: Given any unknown command (e.g. `foo bar`), the engine does not crash and does not emit output (ignores gracefully).

---

## AC-15: TUI Rendering and Game Loop

**Story**: US-24–US-28

**AC-15-01**: When the TUI launches, it renders an 8x8 ASCII board with:
- Correct piece symbols (Unicode or ASCII)
- Rank numbers 1–8 and file letters a–h
- "White to move" indicator

**AC-15-02**: When the player enters a legal move and presses Enter, then:
- The board updates to reflect the move
- "Engine thinking..." is displayed
- The engine's response move is applied
- The updated board is re-rendered

**AC-15-03**: When the player enters an illegal move, then:
- The error message includes the move string
- The board is unchanged
- The move prompt is shown again

**AC-15-04**: When a move delivers check, then "Check!" is displayed after the board render.

**AC-15-05**: When checkmate occurs, then:
- "Checkmate! {winner} wins." is displayed
- No move prompt is shown
- The player is offered to save PGN

---

## AC-16: SSR GUI

**Story**: US-29–US-33

**AC-16-01**: `GET /` returns HTTP 200 with an HTML page containing a chess board in starting position and a "New Game" element.

**AC-16-02**: `POST /game/new` returns HTTP 302 redirecting to `/game/{id}`. `GET /game/{id}` returns HTTP 200 with the board in starting position.

**AC-16-03**: `POST /game/{id}/move?from=e2&to=e4` when the move is legal:
- Returns HTTP 302 to `GET /game/{id}`
- The rendered board shows the pawn on e4
- The engine's response move is already reflected

**AC-16-04**: `POST /game/{id}/move?from=e2&to=e5` when the move is illegal:
- Returns HTTP 422
- The response body contains a human-readable error message
- The board state is unchanged on the next `GET /game/{id}`

**AC-16-05**: When the game is over, `GET /game/{id}` displays the result text and a "New Game" link, with no move input elements present.
