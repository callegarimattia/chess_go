# language: en
Feature: Milestone 4 — Server-Side Rendered GUI
  As Priya the web player
  I want to play chess in a browser without installing any software
  So that I can enjoy a game by simply opening a URL

  # Stories: US-29 through US-33
  # Acceptance Criteria: AC-16
  #
  # All scenarios are tagged @skip.
  # Enable one at a time, implement, commit, then enable the next.
  #
  # Implementation note: SSR scenarios use net/http/httptest.NewServer.
  # They drive web.NewServer(engineFn) through its ServeHTTP port.
  # No real network is opened. The engine is injected as a stub for determinism.

  # ─── Home Page (US-29, AC-16-01) ──────────────────────────────────────────

  @skip
  Scenario: Web player opens the home page and sees the starting board
    Given the SSR server is running
    When a browser requests "GET /"
    Then the response status is 200
    And the response body is an HTML page
    And the HTML contains a chess board in the starting position
    And the HTML contains a "New Game" element

  @skip
  Scenario: Web player sees the home page without JavaScript
    Given the SSR server is running
    When a browser requests "GET /" with JavaScript disabled
    Then the response status is 200
    And the page displays a playable board using only HTML and CSS

  @skip
  Scenario: Web player receives a fast response from the home page
    Given the SSR server is running
    When a browser requests "GET /"
    Then the response is delivered in under 200 milliseconds

  # ─── Create Game Session (US-30, AC-16-02) ────────────────────────────────

  @skip
  Scenario: Web player clicks New Game and is redirected to a fresh game
    Given the SSR server is running
    When a browser sends "POST /game/new"
    Then the response status is 302
    And the Location header contains "/game/" followed by a session identifier

  @skip
  Scenario: Web player follows the new game redirect and sees the starting board
    Given the SSR server is running
    When a browser sends "POST /game/new" and follows the redirect
    Then the response status is 200
    And the board shows the starting position with all 32 pieces

  @skip
  Scenario: Web player can create multiple independent game sessions
    Given the SSR server is running
    When I create two separate game sessions
    Then each session has a unique identifier
    And each session starts at the starting position independently

  @skip
  Scenario: Web player requesting a non-existent session receives a 404 response
    Given the SSR server is running
    When a browser requests "GET /game/nonexistent-session-id"
    Then the response status is 404

  # ─── Make a Move via POST (US-31, AC-16-03) ───────────────────────────────

  @skip
  Scenario: Web player makes a valid move and sees the updated board after the engine responds
    Given the SSR server is running with a deterministic engine stub
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e4"
    Then the response status is 302
    And the Location header is "/game/{id}"
    When the browser follows the redirect to "GET /game/{id}"
    Then the board shows the White pawn on e4
    And the engine stub's response move is already reflected on the board

  @skip
  Scenario: Web player makes a move and the engine responds within the same HTTP request
    Given the SSR server is running with a stub engine that records call timing
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e4"
    Then the engine was called before the 302 redirect was returned

  @skip
  Scenario: Web player sees the full page re-render on the GET after a move
    Given the SSR server is running with a deterministic engine stub
    And a game session exists at the starting position
    When the player makes the move "e2" to "e4" and follows the redirect
    Then the response is a complete HTML page, not a partial fragment
    And the page shows whose turn it is to move

  @skip
  Scenario: Web player receives a fast response when making a move
    Given the SSR server is running with a stub engine capped at 150 milliseconds
    And a game session exists at the starting position
    When a browser sends a valid move POST
    Then the complete request-redirect-render cycle completes in under 200 milliseconds

  # ─── Invalid Move (US-31, AC-16-04) ───────────────────────────────────────

  @skip
  Scenario: Web player submits an illegal move and receives a 422 response
    Given the SSR server is running
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=e2&to=e5"
    Then the response status is 422
    And the response body contains a human-readable error message

  @skip
  Scenario: Web player submits an illegal move and the board state is unchanged
    Given the SSR server is running
    And a game session exists at the starting position
    When a browser sends an illegal move POST and then requests "GET /game/{id}"
    Then the board still shows the starting position

  @skip
  Scenario: Web player submits a move with a missing from parameter and receives a 422 response
    Given the SSR server is running
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?to=e4" without a from parameter
    Then the response status is 422

  @skip
  Scenario: Web player submits a move with an out-of-range square and receives a 422 response
    Given the SSR server is running
    And a game session exists at the starting position
    When a browser sends "POST /game/{id}/move?from=z9&to=e4"
    Then the response status is 422

  @skip
  Scenario: Web player tries to move an opponent's piece and receives a 422 response
    Given the SSR server is running
    And a game session exists at the starting position with White to move
    When a browser sends "POST /game/{id}/move?from=e7&to=e5"
    Then the response status is 422

  # ─── Game Over Page (US-33, AC-16-05) ─────────────────────────────────────

  @skip
  Scenario: Web player sees the game result when the game ends by checkmate
    Given the SSR server is running
    And a game session exists in a position where the next move is checkmate
    When the checkmating move is applied and the player requests "GET /game/{id}"
    Then the response status is 200
    And the HTML contains a result message with the word "Checkmate"
    And the HTML contains a "New Game" link
    And no move input form is present in the HTML

  @skip
  Scenario: Web player sees the game result when the game ends by stalemate
    Given the SSR server is running
    And a game session exists in a stalemate position
    When the player requests "GET /game/{id}"
    Then the HTML contains the word "Stalemate"
    And the HTML contains a "New Game" link
    And no move input form is present in the HTML

  @skip
  Scenario: Web player sees a draw result with the specific draw reason
    Given the SSR server is running
    And a game session exists in a position drawn by the fifty-move rule
    When the player requests "GET /game/{id}"
    Then the HTML contains a draw message
    And the HTML contains a "New Game" link

  @skip
  Scenario: Web player cannot submit a move to a completed game session
    Given the SSR server is running
    And a game session exists where the game has ended by checkmate
    When a browser sends a move POST to that completed session
    Then the response status is 422
    And the board state is not changed

  # ─── SSR Board Content (US-29, FR-05-10) ──────────────────────────────────

  @skip
  Scenario: Web player sees the board rendered as an HTML table
    Given the SSR server is running
    And a game session exists at the starting position
    When a browser requests "GET /game/{id}"
    Then the HTML contains a table element representing the chess board
    And the table has 8 rows and 8 columns
    And piece characters appear in the appropriate cells

  @skip
  Scenario: Web player sees no JavaScript in the HTML response
    Given the SSR server is running
    When a browser requests "GET /"
    Then the HTML does not contain any script elements
    And the board is usable without JavaScript
