# language: en
Feature: Integration Checkpoints — Cross-Layer Scenarios
  As a player or integrator
  I want the full pipeline from input to output to behave correctly end-to-end
  So that the system works as a coherent whole, not just as isolated parts

  # These scenarios deliberately cross layer boundaries to validate that the
  # component contracts are honoured in combination.
  # Each scenario targets a specific integration point documented in the architecture.
  #
  # All scenarios are tagged @skip.
  # Enable one at a time, implement, commit, then enable the next.

  # ─── Chess Package -> Engine Integration ──────────────────────────────────

  @skip
  Scenario: Engine receives a GameState from the chess package and returns a legal move
    Given the chess package produces a GameState from the starting FEN
    When I pass that GameState to the engine Search function
    Then the engine returns a SearchResult whose BestMove is in the LegalMoves of the GameState

  @skip
  Scenario: Engine applies its bestmove through the chess package and produces a valid next state
    Given the starting position GameState
    When the engine selects a bestmove and that move is applied via chess.Game.Apply
    Then the resulting Game has the active color switched to Black
    And the resulting Game's FEN round-trips correctly through ToFEN and NewGameFromFEN

  @skip
  Scenario: Engine search is deterministic for a given position when seeded consistently
    Given the starting position and a fixed random seed
    When I call Search twice with the same movetime
    Then both calls return the same BestMove

  @skip
  Scenario: Engine handles a position produced by a sequence of chess package Apply calls
    Given the starting position
    When I apply "e2e4" then "e7e5" then "g1f3" through the chess package
    And I pass the resulting GameState to the engine Search
    Then the engine returns a legal bestmove for the resulting position
    And no error is returned

  # ─── Chess Package -> TUI Integration ─────────────────────────────────────

  @skip
  Scenario: TUI renderer correctly displays a position produced by the chess package
    Given the chess package produces a GameState after applying "e2e4" to the starting position
    When the TUI renderer writes that GameState to an output buffer
    Then the output shows a pawn on the e4 square
    And the e2 square shows as empty
    And the output indicates it is Black to move

  @skip
  Scenario: TUI game loop calls chess.Game.Apply for every player-entered move
    Given a TUI session is started with a recorder engine stub
    When the player enters the moves "e2e4" and "g1f3"
    Then the chess package Apply was called exactly twice for player moves
    And each call received the correct predecessor GameState

  @skip
  Scenario: TUI correctly displays check when the chess package reports InCheck as true
    Given a TUI session at a position where the player's move will put the opponent in check
    When the player applies the checking move
    Then the TUI output contains a check notification
    And the chess package InCheck was called after the move was applied

  # ─── Chess Package -> SSR Web Integration ─────────────────────────────────

  @skip
  Scenario: SSR handler passes the correct GameState to the engine after a player move
    Given the SSR server is running with an engine stub that records its input
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e4"
    Then the engine stub received the GameState after e2e4 was applied
    And the engine stub's received GameState has Black to move
    And the engine stub's received GameState has the en passant square e3

  @skip
  Scenario: SSR session store returns the same GameState that was stored after a move
    Given the SSR server is running with a deterministic engine stub
    And a game session exists at the starting position
    When the player makes the move "e2e4" and the engine responds with "e7e5"
    And the player makes the move "g1f3" and the engine responds with "b8c6"
    Then "GET /game/{id}" shows the board position after all four moves are applied

  @skip
  Scenario: SSR concurrent requests to different sessions do not interfere
    Given the SSR server is running
    And two separate game sessions exist: session A and session B
    When session A makes the move "e2e4" concurrently with session B making the move "d2d4"
    Then "GET /game/A" shows only the e4 pawn advance
    And "GET /game/B" shows only the d4 pawn advance

  # ─── Engine -> TUI Full Loop Integration ──────────────────────────────────

  @skip
  Scenario: Full TUI game loop: player moves, engine responds, result detected
    Given a TUI session with the engine playing "e7e5" in response to any move
    When the player plays "e2e4"
    Then the board shows both "e4" and "e5" pawns
    And it is White to move again
    And the TUI loop is ready for the next player input

  @skip
  Scenario: Full TUI game loop ends when the engine delivers checkmate
    Given a TUI session at a position where the engine's next move is checkmate
    When the engine move is applied
    Then the TUI displays the checkmate result message
    And the game loop exits without prompting for another move
    And the TUI offers to save the PGN

  # ─── Engine -> SSR Full Request Cycle ─────────────────────────────────────

  @skip
  Scenario: Full SSR request cycle: player move applied, engine responds, board updated in one POST
    Given the SSR server is running with an engine that always plays "g8f6"
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e4"
    And the browser follows the redirect to "GET /game/{id}"
    Then the board shows the White pawn on e4
    And the board shows the Black knight on f6
    And it is White to move

  @skip
  Scenario: SSR server handles engine move that delivers check to the player
    Given the SSR server is running with an engine that always delivers check
    And a game session exists at the position where the engine move gives check
    When a browser sends a valid player move POST and follows the redirect
    Then the HTML contains a check notification visible to the player

  # ─── Error Path Integration ────────────────────────────────────────────────

  @skip
  Scenario: Illegal move is rejected at the chess package boundary before reaching the engine
    Given the SSR server is running with an engine stub that records whether it was called
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e5"
    Then the response status is 422
    And the engine stub was never called

  @skip
  Scenario: Invalid FEN in a UCI position command is rejected before the engine searches
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position fen not-a-valid-fen"
    And I send "go movetime 200"
    Then the engine does not crash
    And the engine either emits bestmove 0000 or skips the search gracefully

  @skip
  Scenario: Engine stop command during active search does not leave the session in a corrupted state
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position startpos"
    And I send "go infinite"
    And I send "stop" after 100 milliseconds
    And I send "position startpos moves e2e4"
    And I send "go movetime 200"
    Then the engine returns a legal bestmove for the position after e2e4
    And no output from the interrupted search contaminates the new search output

  # ─── GameState Immutability Across Layers ─────────────────────────────────

  @skip
  Scenario: The original game state is unchanged after the engine performs a search
    Given the chess package produces a GameState from "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
    When I pass that GameState to the engine Search function and wait for the result
    Then the original GameState still has Black to move
    And the original GameState's FEN is still "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

  @skip
  Scenario: Multiple goroutines can hold references to the same GameState safely
    Given a GameState from the starting position
    When I launch 10 goroutines that each call LegalMoves on the same GameState simultaneously
    Then all 10 goroutines return exactly 20 moves each
    And the go test race detector reports no data races
