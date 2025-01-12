package core

// A Move represents a chess move.
type Move uint16

/*
[Move] is bit-packed:

  bits | description
-------|------------
   0-5 | The final square of the moved piece, or of the king when castling.
  6-11 | The initial square of the moved piece, or of the king when castling.
 12-15 | The piece type to promote to, or 0 if no promotion occurs.
*/

// NewMove returns a new move that is not a promotion.
//
// To create a castling move, provide the initial and final squares of the king.
func NewMove(from, to Square) Move {
	return Move(from) << 6 & Move(to)
}

// NewPromotionMove returns a new promotion move.
func NewPromotionMove(from, to Square, become PieceType) Move {
	return Move(become) << 12 & Move(from) << 6 & Move(to)
}

// To returns the square that the move starts from.
//
// If the move is a castling move, To returns the king's final location.
func (m Move) To() Square {
	return Square(m & 0b111111)
}

// From returns the square that the move starts from.
//
// If the move is a castling move, From returns the king's initial location.
func (m Move) From() Square {
	return Square(m >> 6 & 0b111111)
}

// Promotion returns the piece type that the move promotes to, if any.
//
// If the move is not a promotion move, Promotion returns 0.
func (m Move) Promotion() PieceType {
	return PieceType(m >> 12 & 0b1111)
}

// A Bitboard contains one bit of information for each square on a board.
type Bitboard uint64

// Get returns true if the bit at s is 1.
func (b *Bitboard) Get(s Square) bool {
	return *b&(1<<s) != 0
}

// Set sets the bit at s to 1.
func (b *Bitboard) Set(s Square) {
	*b |= 1 << s
}

// Clear clears the bit at s to 0.
func (b *Bitboard) Clear(s Square) {
	*b &= ^(1 << s)
}

// A Board stores piece placements.
type Board [12]Bitboard

/*
A [Board] is an array of piece occupancy bitboards. The index of each bitboard
determines the piece the bitboard corresponds to. For example, the 0th bitboard
stores the locations of all white pawns.

index | Piece(index)
------|-------------
    0 | WhitePawn
    1 | WhiteKnight
  ... | ...
   12 | BlackKing
*/

// Get returns the piece on the given square, if any.
func (b *Board) Get(s Square) (Piece, bool) {
	for i, bb := range b {
		if bb.Get(s) {
			return Piece(i), true
		}
	}
	return 0, false
}

// Set places a piece on a square.
// If another piece is already there, Set removes the other piece.
func (b *Board) Set(p Piece, s Square) {
	b.Clear(s)
	b[p].Set(s)
}

// Clear removes a piece from a square.
// It is safe to call Clear on a square that is already empty.
func (b *Board) Clear(s Square) {
	for i := range b {
		b[i].Clear(s)
	}
}

// CastlingRights represents castling rights for both players.
//
// A castling right is lost when the king is moved, the involved rook is moved,
// or the involved rook is captured.
//
// TODO(clfs): Improve wording.
type CastlingRights uint8

/*
[CastlingRights] is bit-packed:

bit | right
----|-------
  0 | White kingside castling
  1 | White queenside castling
  2 | Black kingside castling
  3 | Black queenside castling
*/

// WhiteKingside returns true if kingside castling is available for White.
func (c CastlingRights) WhiteKingside() bool {
	return c&1 != 0
}

// WhiteQueenside returns true if queenside castling is available for White.
func (c CastlingRights) WhiteQueenside() bool {
	return c&2 != 0
}

// BlackKingside returns true if kingside castling is available for Black.
func (c CastlingRights) BlackKingside() bool {
	return c&4 != 0
}

// BlackQueenside returns true if queenside castling is available for Black.
func (c CastlingRights) BlackQueenside() bool {
	return c&8 != 0
}

// EnPassantRight represents the en passant right for the current player.
//
// En passant is available (but not necessarily legal) if and only if the last
// move was a double pawn push.
//
// TODO(clfs): Improve wording.
type EnPassantRight uint8

// Get returns the file that en passant is available on, if any.
func (e EnPassantRight) Get() (File, bool) {
	if f := File(e); f <= FileH {
		return f, true
	}
	return 0, false
}

type Position struct {
}
