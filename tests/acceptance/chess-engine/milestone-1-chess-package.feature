# language: en
Feature: Milestone 1 — Chess Logic Package
  As Marco the library integrator
  I want a correct, complete chess rules implementation
  So that I can build move selection and game management logic on top of it

  # Stories: US-01 through US-12
  # Acceptance Criteria: AC-01 through AC-11
  #
  # All scenarios are tagged @skip.
  # Enable one at a time, implement, commit, then enable the next.

  # ─── FEN Parsing (US-01, AC-01) ───────────────────────────────────────────

  @skip
  Scenario: Library consumer parses the starting position from a FEN string
    Given the FEN string "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    When I call NewGameFromFEN with that string
    Then I receive a valid Game with White to move
    And the castling rights are set to all four sides
    And the half-move clock is 0
    And the full-move number is 1
    And all 32 pieces are on their starting squares

  @skip
  Scenario: Library consumer receives a typed error for a malformed FEN string
    Given the FEN string "not-a-valid-fen"
    When I call NewGameFromFEN with that string
    Then I receive an ErrInvalidFEN error
    And no Game is returned
    And no panic occurs

  @skip
  Scenario: Library consumer receives a typed error for a FEN with wrong field count
    Given the FEN string "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq"
    When I call NewGameFromFEN with that string
    Then I receive an ErrInvalidFEN error

  @skip
  Scenario: Library consumer receives a typed error for a FEN with invalid piece characters
    Given the FEN string "rnbqkxnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    When I call NewGameFromFEN with that string
    Then I receive an ErrInvalidFEN error

  @skip
  Scenario: Game state serialises back to the same FEN string it was loaded from
    Given a valid FEN string "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
    When I load the game and immediately call ToFEN
    Then the output FEN string matches the original FEN string

  @skip
  Scenario: FEN round-trip preserves en passant square after a pawn double advance
    Given the starting position
    When I apply the move "e2e4"
    And I call ToFEN on the resulting game
    Then the FEN contains the en passant square "e3"
    And loading that FEN produces an equal game state

  # ─── Legal Move Generation (US-02, AC-02) ─────────────────────────────────

  @skip
  Scenario: Library consumer receives exactly 20 legal moves from the starting position
    Given the starting position
    When I call LegalMoves
    Then exactly 20 moves are returned
    And all returned moves are pawn or knight moves

  @skip
  Scenario: Library consumer receives no illegal moves in any generated move list
    Given the position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I call LegalMoves
    Then every returned move leaves the king not in check after it is applied
    And no returned move places the active king on an attacked square

  @skip
  Scenario: Library consumer receives an empty move list from a checkmate position
    Given the Fool's Mate position "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
    When I call LegalMoves
    Then zero moves are returned

  @skip
  Scenario: Library consumer receives an empty move list from a stalemate position
    Given the stalemate position "k7/8/1Q6/8/8/8/8/7K w - - 0 1"
    When I call LegalMoves
    Then zero moves are returned

  # ─── Apply Move / Immutability (US-03, AC-03) ─────────────────────────────

  @skip
  Scenario: Library consumer applies a pawn double advance and receives the correct new state
    Given the starting position
    When I apply the move "e2e4"
    Then a new game is returned with the pawn on e4
    And it is Black to move in the new game
    And the en passant square is e3 in the new game
    And the original game still has White to move
    And the pawn is still on e2 in the original game

  @skip
  Scenario: Library consumer receives a typed error when applying a move not in the legal move list
    Given the starting position
    When I attempt to apply the move "e2e5"
    Then I receive an ErrIllegalMove error
    And the game state is unchanged

  @skip
  Scenario: Library consumer applies a sequence of moves and each new state is independent
    Given the starting position
    When I apply "e2e4" to get game A
    And I apply "e7e5" to game A to get game B
    And I apply "g1f3" to game B to get game C
    Then game A still has Black to move
    And game B still has White to move
    And game C has the knight on f3 and Black to move

  # ─── Check and Checkmate Detection (US-04, AC-04) ─────────────────────────

  @skip
  Scenario: Library consumer detects checkmate in the Fool's Mate position
    Given the Fool's Mate position "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
    When I call Result on the game
    Then the result is Checkmate
    And the losing side is White

  @skip
  Scenario: Library consumer detects that the active king is in check
    Given the position "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
    When I call InCheck
    Then InCheck returns true

  @skip
  Scenario: Library consumer detects that the active king is not in check
    Given the starting position
    When I call InCheck
    Then InCheck returns false

  @skip
  Scenario: Library consumer detects that a game in progress has no terminal result
    Given the starting position
    When I call Result
    Then the result is InProgress

  # ─── Draw Detection (US-05, AC-05) ────────────────────────────────────────

  @skip
  Scenario: Library consumer detects a draw by the fifty-move rule
    Given a position where the half-move clock is 100
    When I call Result
    Then the result is DrawFiftyMove

  @skip
  Scenario: Library consumer detects a draw by threefold repetition
    Given a game where the knight has shuffled back and forth until the same position has occurred three times
    When I call Result
    Then the result is DrawThreefoldRepetition

  @skip
  Scenario: Library consumer detects a draw by insufficient material with kings only
    Given the position "8/8/8/8/8/8/8/K6k w - - 0 1"
    When I call Result
    Then the result is DrawInsufficientMaterial

  @skip
  Scenario: Library consumer detects stalemate when there are no legal moves and no check
    Given the stalemate position "k7/8/1Q6/8/8/8/8/7K b - - 0 1"
    When I call Result
    Then the result is Stalemate

  # ─── Castling (US-06, AC-06) ──────────────────────────────────────────────

  @skip
  Scenario: Library consumer applies kingside castling for White
    Given a position where White can castle kingside "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"
    When I apply the move "e1g1"
    Then the White king is on g1
    And the White kingside rook is on f1
    And the castling rights no longer include White kingside

  @skip
  Scenario: Library consumer applies queenside castling for White
    Given a position where White can castle queenside "r3kbnr/ppp1pppp/2nqb3/3p4/3P4/2NQB3/PPP1PPPP/R3KBNR w KQkq - 4 5"
    When I apply the move "e1c1"
    Then the White king is on c1
    And the White queenside rook is on d1
    And the castling rights no longer include White queenside

  @skip
  Scenario: Library consumer receives an error when castling through an attacked square
    Given a position where the f1 square is attacked by a Black bishop "rnbqk2r/pppp1ppp/5n2/4p3/1b2P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 4 4"
    When I attempt the move "e1g1"
    Then I receive an ErrIllegalMove error

  @skip
  Scenario: Library consumer receives an error when trying to castle while in check
    Given a position where White is in check with castling rights still set
    When I attempt to castle
    Then I receive an ErrIllegalMove error

  @skip
  Scenario: Library consumer sees castling absent from legal moves after king has moved
    Given a position after the White king has moved and returned
    When I call LegalMoves
    Then no castling moves are present

  # ─── En Passant (US-07, AC-07) ────────────────────────────────────────────

  @skip
  Scenario: Library consumer captures en passant and the captured pawn is removed
    Given a position where White has a pawn on d5 and Black just advanced a pawn to e5 "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq d6 0 3"
    When I apply the move "e5d6"
    Then the White pawn is on d6
    And the Black pawn on d5 is removed from the board
    And the en passant square in the new game is empty

  @skip
  Scenario: Library consumer sees en passant absent from legal moves when the capture would expose the king
    Given a position where capturing en passant would leave the king in check on the same rank
    When I call LegalMoves
    Then the en passant capture is not present in the move list

  # ─── Pawn Promotion (US-08, AC-08) ────────────────────────────────────────

  @skip
  Scenario: Library consumer promotes a White pawn to a queen
    Given a position with a White pawn on e7 "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
    When I apply the move "e7e8q"
    Then a White queen is on e8
    And no pawn is on e7 or e8

  @skip
  Scenario: Library consumer sees four promotion moves for every reachable promotion square
    Given a position with a White pawn on e7 "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
    When I call LegalMoves
    Then four promotion moves to e8 are present: queen, rook, bishop, and knight

  @skip
  Scenario: Library consumer receives an error when applying a promotion move without specifying the piece
    Given a position with a White pawn on e7 "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
    When I attempt to apply the move "e7e8" without a promotion piece
    Then I receive an ErrIllegalMove error

  # ─── Notation Export (US-09, US-10, AC-09) ────────────────────────────────

  @skip
  Scenario: Library consumer receives correct UCI notation for a regular pawn move
    Given the starting position
    When I apply the move "e2e4" and read the UCI string of that move
    Then the UCI string is "e2e4"

  @skip
  Scenario: Library consumer receives correct UCI notation for a promotion move
    Given a position with a White pawn on e7 "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
    When I apply the move "e7e8q" and read the UCI string of that move
    Then the UCI string is "e7e8q"

  @skip
  Scenario: Library consumer receives correct UCI notation for kingside castling
    Given a position where White can castle kingside "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"
    When I apply the move "e1g1" and read the UCI string of that move
    Then the UCI string is "e1g1"

  @skip
  Scenario: Library consumer receives correct SAN notation for a pawn move
    Given the starting position
    When I apply the move "e2e4"
    Then the SAN string of that move is "e4"

  @skip
  Scenario: Library consumer receives correct SAN notation for a knight move
    Given the starting position
    When I apply the move "g1f3"
    Then the SAN string of that move is "Nf3"

  @skip
  Scenario: Library consumer receives correct SAN notation for kingside castling
    Given a position where White can castle kingside "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"
    When I apply the move "e1g1"
    Then the SAN string of that move is "O-O"

  @skip
  Scenario: Library consumer receives correct SAN notation for queenside castling
    Given a position where White can castle queenside "r3kbnr/ppp1pppp/2nqb3/3p4/3P4/2NQB3/PPP1PPPP/R3KBNR w KQkq - 4 5"
    When I apply the move "e1c1"
    Then the SAN string of that move is "O-O-O"

  @skip
  Scenario: Library consumer receives correct SAN notation for a promotion move
    Given a position with a White pawn on e7 "4k3/4P3/8/8/8/8/8/4K3 w - - 0 1"
    When I apply the move "e7e8q"
    Then the SAN string of that move is "e8=Q"

  @skip
  Scenario: Library consumer receives correct SAN notation for a checkmate move
    Given a position one move from Fool's Mate "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq g3 0 2"
    When I apply the move "d8h4"
    Then the SAN string of that move is "Qh4#"

  # ─── PGN Export (US-11, AC-10) ────────────────────────────────────────────

  @skip
  Scenario: Library consumer exports a complete game as a valid PGN string
    Given a completed game of Fool's Mate
    When I call ToPGN
    Then the PGN contains the seven required tag pairs: Event Site Date Round White Black Result
    And the move text is in SAN format with move numbers
    And the result token at the end is "0-1"

  @skip
  Scenario: Library consumer exports an in-progress game and the PGN has a star result token
    Given the starting position with three moves played
    When I call ToPGN
    Then the PGN result token is "*"

  # ─── Perft Validation (US-12, AC-11) ──────────────────────────────────────

  @skip
  Scenario: Move generator produces exactly 20 nodes at depth 1 from the starting position
    Given the starting position
    When I run perft at depth 1
    Then the node count is 20

  @skip
  Scenario: Move generator produces exactly 400 nodes at depth 2 from the starting position
    Given the starting position
    When I run perft at depth 2
    Then the node count is 400

  @skip
  Scenario: Move generator produces exactly 8902 nodes at depth 3 from the starting position
    Given the starting position
    When I run perft at depth 3
    Then the node count is 8902

  @skip
  Scenario: Move generator produces exactly 197281 nodes at depth 4 from the starting position
    Given the starting position
    When I run perft at depth 4
    Then the node count is 197281

  @skip @slow
  Scenario: Move generator produces exactly 4865609 nodes at depth 5 from the starting position
    Given the starting position
    When I run perft at depth 5
    Then the node count is 4865609

  @skip
  Scenario: Move generator produces correct node counts from the Kiwipete position at depth 1
    Given the Kiwipete position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I run perft at depth 1
    Then the node count is 48

  @skip
  Scenario: Move generator produces correct node counts from the Kiwipete position at depth 2
    Given the Kiwipete position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I run perft at depth 2
    Then the node count is 2039

  @skip
  Scenario: Move generator produces correct node counts from the Kiwipete position at depth 3
    Given the Kiwipete position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I run perft at depth 3
    Then the node count is 97862

  @skip
  Scenario: Move generator produces correct node counts from the Kiwipete position at depth 4
    Given the Kiwipete position "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
    When I run perft at depth 4
    Then the node count is 4085603
