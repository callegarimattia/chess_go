# language: en
Feature: Chess Engine in Go
  As a player, library integrator, or engine developer
  I want a correct and usable chess engine
  So that I can play, integrate, or tune a chess AI in Go

  Background:
    Given the chess engine is compiled and available as a binary and Go package

  # ─────────────────────────────────────────────
  # WALKING SKELETON (Feature 0)
  # ─────────────────────────────────────────────

  Scenario: Walking skeleton — engine plays a legal move from starting position
    Given the starting chess position
    When the engine selects and applies a move
    Then the updated board is displayed in the terminal
    And the move selected is a legal chess move
    And no panic or error occurs

  # ─────────────────────────────────────────────
  # CHESS LOGIC PACKAGE — P1: Library Integrator
  # ─────────────────────────────────────────────

  Scenario: Parse a valid FEN string into a GameState
    Given the FEN string "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    When I call NewGameFromFEN with that string
    Then I receive a valid GameState with White to move
    And castling rights are set to KQkq
    And the half-move clock is 0

  Scenario: Reject a malformed FEN string
    Given the FEN string "not-a-valid-fen"
    When I call NewGameFromFEN with that string
    Then I receive an ErrInvalidFEN error
    And no GameState is returned

  Scenario: Generate legal moves from the starting position
    Given the starting chess position
    When I call LegalMoves
    Then exactly 20 moves are returned
    And all moves are pawn or knight moves

  Scenario: Apply a legal move and update game state
    Given the starting chess position
    When I apply the move "e2e4"
    Then the pawn is on e4
    And it is Black to move
    And the en passant square is e3

  Scenario: Reject an illegal move with a typed error
    Given the starting chess position
    When I attempt to apply the move "e2e5"
    Then I receive an ErrIllegalMove error
    And the game state is unchanged

  Scenario: Detect checkmate
    Given the position "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
    When I call GameResult
    Then the result is "checkmate"
    And the winner is Black

  Scenario: Detect stalemate
    Given a position where White has no legal moves and is not in check
    When I call GameResult
    Then the result is "stalemate"
    And the result is a draw

  Scenario: Detect draw by fifty-move rule
    Given a position where the half-move clock has reached 100
    When I call GameResult
    Then the result is "draw-fifty-move"

  Scenario: Detect draw by threefold repetition
    Given a game where the same position has occurred three times
    When I call GameResult
    Then the result is "draw-threefold-repetition"

  Scenario: Execute kingside castling
    Given a position where White can castle kingside
    When I apply the move "e1g1"
    Then the White king is on g1
    And the White rook is on f1
    And castling rights no longer include White kingside

  Scenario: Reject castling through check
    Given a position where the castling path passes through an attacked square
    When I attempt to castle
    Then I receive an ErrIllegalMove error

  Scenario: Execute en passant capture
    Given a position where en passant is available on e6
    When I apply the move "d5e6"
    Then the Black pawn on e5 is removed
    And the White pawn is on e6

  Scenario: Execute pawn promotion
    Given a position where a White pawn is on e7
    When I apply the move "e7e8q"
    Then a White queen is on e8
    And the pawn is removed

  # ─────────────────────────────────────────────
  # ENGINE — P3: Engine Developer
  # ─────────────────────────────────────────────

  Scenario: Engine returns a legal bestmove within time limit
    Given the starting chess position
    When I call SearchBestMove with movetime 1000ms
    Then a bestmove in UCI format is returned within 1050ms
    And the move is legal in the given position

  Scenario: Engine emits UCI info lines before bestmove
    Given any chess position
    When I call SearchBestMove
    Then at least one "info depth" line is emitted before bestmove
    And each info line contains depth, score, nodes, and pv fields

  Scenario: Engine respects stop command
    Given the engine is searching with no time limit
    When I send the "stop" command
    Then the engine emits bestmove within 100ms of receiving stop
    And no further output is produced

  # ─────────────────────────────────────────────
  # UCI PROTOCOL — P3: Engine Developer
  # ─────────────────────────────────────────────

  Scenario: UCI handshake completes correctly
    When I send "uci"
    Then the engine responds with "id name" and "id author" lines
    And the engine responds with "uciok"

  Scenario: Engine responds to isready
    Given the UCI handshake is complete
    When I send "isready"
    Then the engine responds with "readyok"

  Scenario: Engine resets state on ucinewgame
    Given the engine has played a game
    When I send "ucinewgame"
    Then the engine's internal state is reset to starting position

  Scenario: Engine processes position with moves
    When I send "position startpos moves e2e4 e7e5"
    And I send "go movetime 500"
    Then the engine returns bestmove
    And the move is legal in the position after e2e4 e7e5

  Scenario: Engine processes position from FEN
    When I send "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
    And I send "go movetime 500"
    Then the engine returns a legal bestmove for Black

  # ─────────────────────────────────────────────
  # TUI — P2: Human TUI Player
  # ─────────────────────────────────────────────

  Scenario: TUI renders starting board correctly
    When I launch the chess-go binary
    Then an ASCII chess board is displayed in the terminal
    And pieces are in their starting positions
    And coordinates a-h and 1-8 are visible
    And "White to move" is displayed

  Scenario: TUI accepts legal move and updates board
    Given the TUI is running with White to move
    When I type "e2e4" and press Enter
    Then the board updates with the pawn on e4
    And "Black to move" is displayed
    And "Engine thinking..." is shown while the engine responds

  Scenario: TUI rejects illegal move with helpful message
    Given the TUI is running with White to move
    When I type "e2e5" and press Enter
    Then "Illegal move: e2e5" is displayed
    And the board is unchanged
    And the move prompt is shown again

  Scenario: TUI displays check notification
    Given a position where the move puts the opponent in check
    When I apply a move that gives check
    Then "Check!" is displayed after the board render

  Scenario: TUI displays game result and prompts for action
    Given a position one move from checkmate
    When White delivers checkmate
    Then "Checkmate! White wins." is displayed
    And "Save game? (y/n)" is prompted

  # ─────────────────────────────────────────────
  # SSR GUI — P4: Casual Web Player
  # ─────────────────────────────────────────────

  Scenario: SSR GUI serves starting board
    When I navigate to "GET /"
    Then I receive an HTML page with a chess board in starting position
    And a "New Game" button is present

  Scenario: SSR GUI creates a new game session
    When I POST to "/game/new"
    Then I am redirected to "/game/{session-id}"
    And the board is rendered in starting position

  Scenario: SSR GUI applies a valid move
    Given a game session at starting position
    When I POST to "/game/{id}/move?from=e2&to=e4"
    Then I am redirected to "GET /game/{id}"
    And the board shows the pawn on e4
    And the engine's response move is already applied

  Scenario: SSR GUI rejects an invalid move
    Given a game session at starting position
    When I POST to "/game/{id}/move?from=e2&to=e5"
    Then I receive an HTTP 422 response
    And the body contains an error message
    And the board state is unchanged

  Scenario: SSR GUI displays game result
    Given a game session where checkmate has occurred
    When I navigate to "GET /game/{id}"
    Then "Checkmate! {winner} wins." is displayed
    And a "New Game" link is present
    And no move input is available
