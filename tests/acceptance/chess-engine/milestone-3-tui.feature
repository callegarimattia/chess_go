# language: en
Feature: Milestone 3 — Terminal User Interface
  As Sofia the TUI player
  I want to play chess in the terminal against the engine with a clear board and feedback
  So that I can enjoy a complete game from launch to result without reading documentation

  # Stories: US-24 through US-28
  # Acceptance Criteria: AC-15
  #
  # All scenarios are tagged @skip.
  # Enable one at a time, implement, commit, then enable the next.
  #
  # Implementation note: TUI scenarios drive the tui package via injected io.Reader and io.Writer.
  # Use tui.NewGame(reader, writer, engineFn).Run() — never spawn a real subprocess.
  # This exercises the driving port (tui.NewGame public API) without bypassing it.

  # ─── Board Render (US-24, AC-15-01) ───────────────────────────────────────

  @skip
  Scenario: TUI player launches the game and sees the starting board
    Given a new TUI game session is started
    When the game loop renders the initial position
    Then the output contains an 8-by-8 grid of piece symbols
    And rank numbers 1 through 8 are visible in the output
    And file letters a through h are visible in the output
    And the text "White to move" appears in the output

  @skip
  Scenario: TUI player sees all starting pieces in their correct positions
    Given a new TUI game session is started
    When the game loop renders the initial position
    Then the White king symbol appears on the e1 position
    And the Black king symbol appears on the e8 position
    And pawns are present on ranks 2 and 7

  @skip
  Scenario: TUI player sees Black to move indicator after the first move
    Given a new TUI game session is started
    When the player inputs the move "e2e4"
    Then the output contains "Black to move" or an equivalent side-to-move indicator

  # ─── Legal Move Acceptance (US-25, AC-15-02) ──────────────────────────────

  @skip
  Scenario: TUI player enters a legal move and sees the updated board
    Given a new TUI game session is started with a predictable engine stub
    When the player inputs the move "e2e4"
    Then the board is redrawn with the pawn on e4
    And the engine thinking indicator is displayed before the engine move is shown

  @skip
  Scenario: TUI player enters a legal move and the engine responds with a move
    Given a new TUI game session is started with a predictable engine stub that always plays "e7e5"
    When the player inputs the move "e2e4"
    Then the board shows the White pawn on e4
    And the board shows the Black pawn on e5
    And it is White to move again

  @skip
  Scenario: TUI player makes several moves in sequence and the game state accumulates correctly
    Given a new TUI game session is started with an engine stub
    When the player inputs the moves "e2e4" then "g1f3" then "f1c4"
    Then the board shows White pieces on e4, f3, and c4

  # ─── Illegal Move Rejection (US-25, AC-15-03) ─────────────────────────────

  @skip
  Scenario: TUI player enters an illegal move and sees a clear error message
    Given a new TUI game session is started
    When the player inputs the move "e2e5"
    Then the output contains the text "Illegal move" and the move string "e2e5"
    And the board is unchanged from the previous render
    And the move prompt is displayed again

  @skip
  Scenario: TUI player enters a move in the wrong format and sees a format error
    Given a new TUI game session is started
    When the player inputs the text "hello"
    Then the output contains a message indicating the move format is invalid
    And the board is unchanged
    And the move prompt is displayed again

  @skip
  Scenario: TUI player enters an empty input and the prompt is shown again
    Given a new TUI game session is started
    When the player presses Enter without typing a move
    Then the move prompt is shown again without changing the board

  @skip
  Scenario: TUI player tries to move a piece that belongs to the opponent
    Given a new TUI game session is started with White to move
    When the player inputs the move "e7e5"
    Then the output contains an illegal move message
    And the board is unchanged

  # ─── Check Notification (US-27, AC-15-04) ─────────────────────────────────

  @skip
  Scenario: TUI player sees a check notification when a move delivers check
    Given a TUI game session at a position one move from delivering check
    When the engine plays a move that puts White in check
    Then the output contains the text "Check!" after the board is rendered

  @skip
  Scenario: TUI player is notified when the player's own move delivers check to the opponent
    Given a TUI game with a position where "e5f7" delivers check
    When the player inputs the move "e5f7"
    Then the output contains "Check!" before the engine responds

  # ─── Game Result Display (US-27, AC-15-05) ────────────────────────────────

  @skip
  Scenario: TUI player sees a checkmate result message and no further move prompt
    Given a TUI game session where the next engine move is checkmate
    When the engine delivers checkmate
    Then the output contains "Checkmate" and the winner's colour
    And no move prompt appears after the result message

  @skip
  Scenario: TUI player sees a stalemate result message when the game is drawn by stalemate
    Given a TUI game session where the next move results in stalemate
    When the move that causes stalemate is played
    Then the output contains "Stalemate"
    And no move prompt appears after the result message

  @skip
  Scenario: TUI player sees a draw-by-fifty-move-rule message
    Given a TUI game where the half-move clock reaches 100
    When Result is checked
    Then the output contains a draw message mentioning the fifty-move rule

  @skip
  Scenario: TUI player sees a draw-by-threefold-repetition message
    Given a TUI game where the same position has occurred three times
    When Result is checked
    Then the output contains a draw message mentioning repetition

  # ─── PGN Export (US-28, AC-15-05) ────────────────────────────────────────

  @skip
  Scenario: TUI player is offered a save prompt when the game ends
    Given a TUI game session where the next move ends the game by checkmate
    When the game ends
    Then the output contains a prompt asking whether to save the PGN
    And the game loop does not exit before the player answers

  @skip
  Scenario: TUI player saves the PGN and the file is written to disk
    Given a TUI game that has ended by checkmate
    When the player answers "y" at the save prompt with a file path
    Then a file is created at the given path
    And the file contents are valid PGN with the seven required tag pairs

  @skip
  Scenario: TUI player declines to save the PGN and the game exits cleanly
    Given a TUI game that has ended by checkmate
    When the player answers "n" at the save prompt
    Then no PGN file is created
    And the process exits with code 0
