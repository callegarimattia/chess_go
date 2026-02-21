# language: en
Feature: Walking Skeleton — End-to-End Pipeline (Feature 0)
  As Daniel the engine developer
  I want the engine to select a random legal move from the starting position and display the updated board in the terminal
  So that I can validate that all architectural layers are connected correctly

  # US-00: End-to-end skeleton
  # AC-00: Walking Skeleton
  #
  # This is the ONLY scenario that is enabled by default.
  # All other scenarios across all feature files are tagged @skip.
  # Enable one at a time following Outside-In TDD.

  @walking_skeleton
  Scenario: Player launches the engine and sees the starting position advance by one legal move
    Given the chess engine is compiled and runnable
    And the starting chess position is loaded
    When the engine selects and applies a move
    Then the updated board is displayed in the terminal
    And the move selected is one of the 20 legal moves from the starting position
    And the process exits with code 0
    And no panic or error occurs

  @walking_skeleton
  Scenario: Walking skeleton pipeline connects all architectural layers without error
    Given the starting FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" is provided
    When the full pipeline runs: FEN parse -> move generation -> random selection -> apply -> render
    Then each pipeline stage produces a valid output
    And the rendered board reflects the applied move
    And the original starting position is unchanged
