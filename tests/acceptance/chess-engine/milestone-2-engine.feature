# language: en
Feature: Milestone 2 — Search Engine
  As Daniel the engine developer
  I want a search engine that finds strong moves within time limits and communicates via UCI
  So that I can connect it to standard chess GUIs and validate its correctness

  # Stories: US-13 through US-23
  # Acceptance Criteria: AC-12, AC-13, AC-14
  #
  # All scenarios are tagged @skip.
  # Enable one at a time, implement, commit, then enable the next.

  # ─── Random Move Engine / Skeleton (US-13) ────────────────────────────────

  @skip
  Scenario: Engine developer gets a legal move from the random move selector
    Given the starting position
    When I call the random move engine with no time limit
    Then a move is returned within 10ms
    And the returned move is in the legal move list for the starting position

  @skip
  Scenario: Engine developer gets a legal move from any non-terminal position
    Given the position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I call the random move engine
    Then the returned move is in the legal move list for that position

  # ─── Alpha-Beta Search (US-14, AC-12) ─────────────────────────────────────

  @skip
  Scenario: Engine developer receives a legal bestmove within the time limit from the starting position
    Given the starting position
    When I call Search with a movetime of 1000 milliseconds
    Then a bestmove is returned within 1050 milliseconds
    And the bestmove is in the legal move list for the starting position

  @skip
  Scenario: Engine developer sees info lines emitted during search
    Given the starting position
    When I call Search with a movetime of 500 milliseconds
    Then at least one info line is emitted before the bestmove
    And each info line contains the fields: depth score nodes nps pv

  @skip
  Scenario: Engine developer sees the search reach at least depth 3 in 100 milliseconds
    Given the starting position
    When I call Search with a movetime of 100 milliseconds
    Then the search reaches at least depth 3
    And a bestmove is returned

  @skip
  Scenario: Engine developer sees the engine find a forced mate in one
    Given the position one move from checkmate "k7/8/1K6/8/8/8/8/R7 w - - 0 1"
    When I call Search with a movetime of 100 milliseconds
    Then the bestmove delivers checkmate

  @skip
  Scenario: Engine developer sees the engine find a forced mate in one in Fool's Mate setup
    Given the position "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq g3 0 2"
    When I call Search with a movetime of 100 milliseconds
    Then the bestmove is "d8h4"

  # ─── Material Evaluation (US-15) ──────────────────────────────────────────

  @skip
  Scenario: Engine developer sees the engine avoid losing a queen for a pawn
    Given a position where White can capture a pawn with the queen but Black can immediately recapture with a pawn
    When I call Search with a movetime of 500 milliseconds
    Then the bestmove does not capture the pawn with the queen on that square

  @skip
  Scenario: Engine developer sees the engine choose a move that captures a free piece
    Given a position where White can capture an undefended Black knight
    When I call Search with a movetime of 200 milliseconds
    Then the bestmove captures the undefended knight

  # ─── Time Management (US-17, AC-13) ───────────────────────────────────────

  @skip
  Scenario: Engine developer sees the bestmove returned within the movetime grace period
    Given the starting position
    When I call Search with a movetime of 500 milliseconds
    Then the bestmove is returned within 550 milliseconds

  @skip
  Scenario: Engine developer sees the engine respect a very short movetime
    Given the starting position
    When I call Search with a movetime of 50 milliseconds
    Then the bestmove is returned within 100 milliseconds
    And the returned move is a legal move

  @skip
  Scenario: Engine developer sees the engine stay within clock allocation over multiple moves
    Given a 60-second time control for both sides
    When the engine plays 30 moves with that time control
    Then no single move consumed more than 10 seconds
    And the engine did not exceed the total allocated time

  # ─── UCI Handshake (US-20, AC-14) ─────────────────────────────────────────

  @skip
  Scenario: Engine developer sends uci and receives the required identification lines
    Given the chess-go binary is started as a subprocess
    When I send the command "uci"
    Then the engine outputs a line matching "id name chess-go"
    And the engine outputs a line matching "id author"
    And the engine outputs the line "uciok"
    And all three responses arrive within 100 milliseconds

  @skip
  Scenario: Engine developer sends isready after uci and receives readyok
    Given the chess-go binary is started and the UCI handshake is complete
    When I send the command "isready"
    Then the engine outputs the line "readyok" within 100 milliseconds

  @skip
  Scenario: Engine developer resets state with ucinewgame between two searches
    Given the chess-go binary is started and a game has been searched
    When I send the command "ucinewgame"
    And I send "position startpos"
    And I send "go movetime 200"
    Then the engine returns a bestmove for the starting position
    And no output from the previous game appears

  # ─── Position Command (US-21) ─────────────────────────────────────────────

  @skip
  Scenario: Engine developer sets up a position from the starting position with moves applied
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position startpos moves e2e4 e7e5"
    And I send "go movetime 500"
    Then the engine returns a bestmove
    And the bestmove is legal in the position after e2e4 e7e5

  @skip
  Scenario: Engine developer sets up a position directly from a FEN string
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
    And I send "go movetime 500"
    Then the engine returns a legal bestmove for Black

  @skip
  Scenario: Engine developer sets up a position from FEN with additional moves
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1 moves e7e5"
    And I send "go movetime 200"
    Then the engine returns a legal bestmove for White

  # ─── Go Command (US-22) ────────────────────────────────────────────────────

  @skip
  Scenario: Engine developer uses go movetime to cap the search duration
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position startpos"
    And I send "go movetime 300"
    Then the engine returns a bestmove within 350 milliseconds

  @skip
  Scenario: Engine developer uses go wtime btime to allocate time from a game clock
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position startpos"
    And I send "go wtime 60000 btime 60000"
    Then the engine returns a bestmove within 5000 milliseconds
    And the bestmove is a legal starting-position move

  # ─── Stop and Quit Commands (US-23, AC-14) ────────────────────────────────

  @skip
  Scenario: Engine developer sends stop during an active search and receives bestmove promptly
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "position startpos"
    And I send "go infinite"
    And I wait 200 milliseconds
    And I send "stop"
    Then the engine outputs a bestmove within 100 milliseconds of receiving stop
    And no further output is produced after the bestmove line

  @skip
  Scenario: Engine developer sends quit and the process exits cleanly
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "quit"
    Then the process exits with code 0 within 500 milliseconds

  @skip
  Scenario: Engine developer sends an unknown command and the engine does not crash
    Given the chess-go binary is started and the UCI handshake is complete
    When I send "foo bar baz"
    Then no output is produced
    And the engine continues to accept commands
    And sending "isready" after the unknown command returns "readyok"

  # ─── Quiescence Search (US-18) ────────────────────────────────────────────

  @skip
  Scenario: Engine developer sees the engine avoid the horizon effect on a known tactical position
    Given a position with a hanging piece that a shallow fixed-depth search would miss "r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 0 4"
    When I call Search with a movetime of 500 milliseconds
    Then the bestmove does not walk into a clearly losing material exchange

  # ─── Move Ordering (US-19) ────────────────────────────────────────────────

  @skip
  Scenario: Engine developer sees captures ordered before quiet moves at every depth
    Given a position with available captures
    When I call Search with a movetime of 200 milliseconds
    Then the search evaluates capture moves before quiet moves at the same depth level
