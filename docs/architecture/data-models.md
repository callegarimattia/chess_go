# Data Models — Chess Engine in Go
**Epic**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved

All models are Go value types (structs) unless noted. No database schema — all state is in-memory.

---

## Core Domain Models (internal/chess)

### Color

```
Color: uint8
  White = 0
  Black = 1
```

### Piece

```
Piece: uint8
  NoPiece = 0
  WhitePawn, WhiteKnight, WhiteBishop, WhiteRook, WhiteQueen, WhiteKing
  BlackPawn, BlackKnight, BlackBishop, BlackRook, BlackQueen, BlackKing
```

Piece encodes both color and type in a single byte. Color extracted by `piece >> 3`. Type extracted by `piece & 0x7`. No separate PieceType enum required.

### Square

```
Square: uint8
  Values: 0–63 (a1=0, b1=1, ..., h8=63)
  NoSquare = 64

  Rank(sq) = sq / 8       (0–7, where 0 = rank 1)
  File(sq) = sq % 8       (0–7, where 0 = file a)
  SquareOf(file, rank) = rank*8 + file
```

### Board

```
Board: [64]Piece
  Index: Square (0–63)
  Board[sq] = Piece at square sq (NoPiece if empty)
```

Board is a value type: copying a Board copies all 64 bytes. No heap allocation.

### CastlingRights

```
CastlingRights: uint8  (bitmask)
  CastleWhiteKingside  = 0b0001
  CastleWhiteQueenside = 0b0010
  CastleBlackKingside  = 0b0100
  CastleBlackQueenside = 0b1000
  NoCastling           = 0b0000
  AllCastling          = 0b1111
```

### GameState (SA-04)

```
GameState:
  Board          Board           // 64-byte piece placement
  ActiveColor    Color           // whose turn to move
  CastlingRights CastlingRights  // current castling availability
  EnPassantSq    Square          // en passant target square (NoSquare if none)
  HalfMoveClock  uint8           // moves since last pawn move or capture (fifty-move rule)
  FullMoveNumber uint16          // starts at 1, increments after Black moves
```

GameState is a value type. Copying GameState is cheap (< 80 bytes total). Immutability is achieved by convention: `Apply()` constructs and returns a new GameState; the original is never mutated.

### Game

```
Game:
  State          GameState       // current position
  History        []GameState     // all positions since game start (for threefold repetition)
  Moves          []Move          // all moves played (for PGN and display)
  ZobristHistory []uint64        // Zobrist hash of each position (fast repetition detection)
```

Game is the main type consumers interact with. It wraps GameState with history for draw detection and export. `Game.Apply()` appends to history and returns a new `Game`.

### Move (SA-02)

```
Move:
  From      Square  // source square
  To        Square  // destination square
  Promotion Piece   // NoPiece if not a promotion; else WhiteQueen/Rook/Bishop/Knight

  UCIString() string          // "e2e4", "e7e8q", "e1g1"
  SANString(g Game) string    // "e4", "Nf3", "O-O", "Bxe5+", "e8=Q#"
  IsCapture(g Game) bool      // true if move captures a piece or en passant
  IsPromotion() bool          // true if Promotion != NoPiece
  IsCastle() bool             // true if king moves 2 squares
```

Move is a value type (3 bytes). No heap allocation in the hot search path.

### GameResult

```
GameResult: uint8
  InProgress             = 0
  WhiteWins (checkmate)  = 1
  BlackWins (checkmate)  = 2
  DrawStalemate          = 3
  DrawFiftyMove          = 4
  DrawThreefoldRepetition = 5
  DrawInsufficientMaterial = 6
```

### Error Types

```
ErrIllegalMove        // returned by Game.Apply() when move not in LegalMoves()
ErrInvalidFEN         // returned by NewGameFromFEN() on malformed input
ErrInvalidMoveFormat  // returned by ParseUCIMove() on malformed UCI string
```

All errors are sentinel values (comparable with `errors.Is`).

---

## Engine Models (internal/engine)

### TimeControl (SA-07)

```
TimeControl:
  MoveTime  time.Duration  // exact time for this move (movetime command); 0 if using wtime/btime
  WTime     time.Duration  // White remaining time
  BTime     time.Duration  // Black remaining time
  WInc      time.Duration  // White increment per move
  BInc      time.Duration  // Black increment per move

  AllocatedTime(color Color) time.Duration  // computed allocation for current search
```

When `MoveTime > 0`, the search uses MoveTime as the hard deadline. Otherwise, a fraction of the relevant side's remaining time is allocated. Typical fraction: `remaining / 30 + increment * 0.8`.

### SearchResult (SA-06)

```
SearchResult:
  BestMove  chess.Move  // best move found in search
  Score     int         // centipawn score (positive = good for active side)
  Depth     int         // depth reached in last completed iteration
  Nodes     int64       // total nodes evaluated
  Elapsed   time.Duration
```

### SearchInfo (emitted as UCI info lines)

```
SearchInfo:
  Depth  int
  Score  int          // centipawns
  Nodes  int64
  NPS    int64        // nodes per second
  PV     []chess.Move // principal variation
```

Each completed depth iteration emits a SearchInfo line to the `info io.Writer` passed to `Search()`.

### KillerMoves

```
KillerMoves: [MaxDepth][2]chess.Move
  // Two killer moves per depth level (quiet moves that caused beta cutoffs)
```

Stored in search context, not passed across package boundaries.

---

## Web Models (internal/web)

### Session

```
Session:
  ID         string      // UUID hex string (32 chars from crypto/rand)
  Game       chess.Game  // current game state (immutable; replaced on each move)
  CreatedAt  time.Time
  LastMoved  time.Time
```

### SessionStore

```
SessionStore:
  mu       sync.RWMutex
  sessions map[string]Session

  Get(id string) (Session, bool)
  Put(id string, s Session)
  New() Session  // generates ID, initializes from starting FEN
```

Sessions are held in memory. No persistence. Server restart loses all sessions (v1 documented limitation).

---

## Shared Artifact Traceability

| Artifact | Model | Package | Key Methods |
|----------|-------|---------|-------------|
| SA-01 FEN string | string | chess | `NewGameFromFEN()`, `Game.ToFEN()` |
| SA-02 Move (UCI) | `chess.Move` | chess | `Move.UCIString()`, `ParseUCIMove()` |
| SA-03 Move (SAN) | string | chess | `Move.SANString(g Game)` |
| SA-04 GameState | `chess.Game` | chess | `Game.Apply()`, `Game.LegalMoves()` |
| SA-05 PGN record | string | chess | `Game.ToPGN()` |
| SA-06 Engine response | `engine.SearchResult` | engine | `Search()` |
| SA-07 Time control | `engine.TimeControl` | engine | `TimeControl.AllocatedTime()` |

---

## Zobrist Hashing (Threefold Repetition)

Zobrist keys are 64-bit pseudo-random values, one per (piece, square) combination plus side-to-move and castling right bits.

```
ZobristTable:
  PieceSquare [12][64]uint64  // 12 piece types × 64 squares
  SideToMove         uint64
  CastlingRights  [16]uint64  // one per CastlingRights bitmask value
  EnPassantFile    [8]uint64  // one per file (a–h)
```

Initialized once at package init with a deterministic seed. Hash of a position:

```
hash = XOR of:
  PieceSquare[piece][sq] for each occupied square
  SideToMove if Black to move
  CastlingRights[castlingRightsBitmask]
  EnPassantFile[epFile] if en passant available
```

The `Game.ZobristHistory` slice stores the hash after each half-move. Threefold repetition is detected when the current hash appears at least twice in history (matching positions, since color and castling are encoded in the hash).

---

## Board Representation Decision

The board uses an **8x8 array** (`[64]Piece`), not bitboards. This decision is documented in ADR-001.

Key properties:
- Piece lookup by square: O(1) array index
- Iteration over all pieces: O(64) scan of array
- Move generation: directional iteration with boundary checks using file/rank arithmetic
- No bitwise operations on 64-bit integers required
- Sufficient for 100k+ NPS target on modern hardware

The 8x8 array is the simplest representation that satisfies v1 requirements. Bitboard migration is a v2 decision after profiling demonstrates the array is the actual bottleneck.
