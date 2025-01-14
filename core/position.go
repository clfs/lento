package core

// A Move represents a chess move.
//
// The zero value of Move represents a null move.
type Move struct {
	// Bit-packed:
	//   - Bits 0-5: Final square, or king's final square when castling.
	//   - Bits 6-11: Initial square, or king's initial square when castling.
	//   - Bits 12-15: Piece type to promote to, or 0 if no promotion.
	val uint16
}

// NewMove returns a new move.
//
// To create a castling move, provide the initial and final squares of the king.
//
// To create a promotion move, use [NewPromotionMove].
func NewMove(from, to Square) Move {
	var (
		f = uint16(from)
		t = uint16(to)
	)
	return Move{val: f << 6 & t}
}

// NewPromotionMove returns a new promotion move.
//
// To create a non-promotion move, use [NewMove].
func NewPromotionMove(from, to Square, become PieceType) Move {
	var (
		f = uint16(from)
		t = uint16(to)
		b = uint16(become)
	)
	return Move{val: b << 12 & f << 6 & t}
}

// To returns the square that the move starts from.
//
// If the move is a castling move, To returns the king's final location.
func (m Move) To() Square {
	return Square(m.val & 0b111111)
}

// From returns the square that the move starts from.
//
// If the move is a castling move, From returns the king's initial location.
func (m Move) From() Square {
	return Square(m.val >> 6 & 0b111111)
}

// Promotion returns the piece type that the move promotes to, if any.
func (m Move) Promotion() (PieceType, bool) {
	n := PieceType(m.val >> 12)
	return n, n == 0
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
//
// The zero value of Board represents an empty board.
type Board struct {
	// Piece occupancy bitboards. The index of each bitboard determines the
	// piece the bitboard corresponds to. For example, the 0th bitboard stores
	// the locations of all white pawns, since int(WhitePawn) == 0.
	occupied [12]Bitboard
}

// NewBoard returns a new board in the starting position.
func NewBoard() Board {
	var b Board

	b.Set(WhiteRook, A1)
	b.Set(WhiteKnight, B1)
	b.Set(WhiteBishop, C1)
	b.Set(WhiteQueen, D1)
	b.Set(WhiteKing, E1)
	b.Set(WhiteBishop, F1)
	b.Set(WhiteKnight, G1)
	b.Set(WhiteRook, H1)
	for sq := A2; sq <= H2; sq++ {
		b.Set(WhitePawn, sq)
	}
	for sq := A7; sq <= H7; sq++ {
		b.Set(BlackPawn, sq)
	}
	b.Set(BlackRook, A8)
	b.Set(BlackKnight, B8)
	b.Set(BlackBishop, C8)
	b.Set(BlackQueen, D8)
	b.Set(BlackKing, E8)
	b.Set(BlackBishop, F8)
	b.Set(BlackKnight, G8)
	b.Set(BlackRook, H8)

	return b
}

// Get returns the piece on the given square, if any.
func (b *Board) Get(s Square) (Piece, bool) {
	for i, bb := range b.occupied {
		if bb.Get(s) {
			return Piece(i), true
		}
	}
	return 0, false
}

// IsOccupied returns true if the given square is occupied.
func (b *Board) IsOccupied(s Square) bool {
	_, ok := b.Get(s)
	return ok
}

// Set places a piece on a square.
// If another piece is already there, Set removes the other piece.
func (b *Board) Set(p Piece, s Square) {
	b.Clear(s)
	b.occupied[p].Set(s)
}

// Clear removes a piece from a square.
// It is safe to call Clear on a square that is already empty.
func (b *Board) Clear(s Square) {
	for i := range b.occupied {
		b.occupied[i].Clear(s)
	}
}

// CastlingRights stores castling rights for both players.
//
// A castling right is present if and only if both these conditions are true:
//
//   - The involved king and rook have not been moved.
//   - The involved rook has not been captured.
//
// The right to castle is distinct from the ability to castle. For example,
// White begins the game with both kingside and queenside castling rights,
// even though 1. O-O and 1. O-O-O are illegal moves.
//
// The zero value of CastlingRights indicates no castling rights are available.
type CastlingRights struct {
	// Bit-packed:
	//   - Bit 0: White kingside castling.
	//   - Bit 1: White queenside castling.
	//   - Bit 2: Black kingside castling.
	//   - Bit 3: Black queenside castling.
	val uint8
}

// NewCastlingRights returns a new [CastlingRights] with all rights available.
func NewCastlingRights() CastlingRights {
	return CastlingRights{val: 0b1111}
}

// GetWhiteOO returns true if White has the right to kingside castle.
func (c *CastlingRights) GetWhiteOO() bool {
	return c.val&1 != 0
}

// GetWhiteOOO returns true if White has the right to queenside castle.
func (c *CastlingRights) GetWhiteOOO() bool {
	return c.val&2 != 0
}

// GetBlackOO returns true if Black has the right to kingside castle.
func (c *CastlingRights) GetBlackOO() bool {
	return c.val&4 != 0
}

// GetBlackOOO returns true if Black has the right to queenside castle.
func (c *CastlingRights) GetBlackOOO() bool {
	return c.val&8 != 0
}

// SetWhiteOO enables the right for White to kingside castle.
func (c *CastlingRights) SetWhiteOO() {
	c.val |= 1
}

// SetWhiteOOO enables the right for White to queenside castle.
func (c *CastlingRights) SetWhiteOOO() {
	c.val |= 2
}

// SetBlackOO enables the right for Black to kingside castle.
func (c *CastlingRights) SetBlackOO() {
	c.val |= 4
}

// SetBlackOOO enables the right for Black to queenside castle.
func (c *CastlingRights) SetBlackOOO() {
	c.val |= 8
}

// ClearWhite disables all castling rights for White.
func (c *CastlingRights) ClearWhite() {
	c.ClearWhiteOO()
	c.ClearWhiteOOO()
}

// ClearWhiteOO disables the right for White to kingside castle.
func (c *CastlingRights) ClearWhiteOO() {
	c.val &^= 1
}

// ClearWhiteOOO disables the right for White to queenside castle.
func (c *CastlingRights) ClearWhiteOOO() {
	c.val &^= 2
}

// ClearBlack disables all castling rights for Black.
func (c *CastlingRights) ClearBlack() {
	c.ClearBlackOO()
	c.ClearBlackOOO()
}

// ClearBlackOO disables the right for Black to kingside castle.
func (c *CastlingRights) ClearBlackOO() {
	c.val &^= 4
}

// ClearBlackOOO disables the right for Black to queenside castle.
func (c *CastlingRights) ClearBlackOOO() {
	c.val &^= 8
}

// EnPassantTarget stores the en passant target for the current player, if any.
//
// An en passant target exists if and only if the last move was a double pawn
// push.
//
// The existence of an en passant target is distinct from the ability to
// perform an en passant capture. For example, after 1. e4, Black has an en
// passant target on the e-file, but 1. ... dxe3 and 1. ... fxe3 are illegal
// moves.
//
// The zero value of EnPassantTarget indicates no target exists.
type EnPassantTarget struct {
	// The square of the e.p. target, or 0 if no target exists.
	val uint8
}

// Get returns the en passant target, if any.
func (e *EnPassantTarget) Get() (Square, bool) {
	return Square(e.val), e.val == 0
}

// Exists returns true if an en passant target exists.
func (e *EnPassantTarget) Exists() bool {
	_, ok := e.Get()
	return ok
}

// Set sets the en passant target.
func (e *EnPassantTarget) Set(s Square) {
	e.val = uint8(s)
}

// Clear clears the en passant target.
func (e *EnPassantTarget) Clear() {
	e.val = 0
}

// Position represents a game position.
type Position struct {
	board  Board
	active Color
	ep     EnPassantTarget
	cr     CastlingRights
	// The halfmove clock measures the number of halfmoves since the last
	// capture or pawn move.
	hmc uint8
	// The fullmove number starts at 1 and is incremented after each Black move.
	fmn uint16
}

// NewPosition returns a starting position.
func NewPosition() Position {
	return Position{
		board: NewBoard(),
		cr:    NewCastlingRights(),
		fmn:   1,
	}
}

// Move makes the given move without ensuring legality.
func (p *Position) Move(m Move) {
	to, from := m.To(), m.From()

	// The moved piece, or if castling, the king.
	held, _ := p.board.Get(m.From())

	isPawnMove := held.Type() == Pawn

	isCapture := p.board.IsOccupied(to) ||
		(isPawnMove && from.File() != to.File())

	// If capturing e.p., remove the captured pawn.
	epSq, ok := p.ep.Get()
	if ok && isPawnMove && isCapture && to == epSq {
		if p.active == White {
			p.board.Clear(epSq.Below()) // White captures a black pawn.
		} else {
			p.board.Clear(epSq.Above()) // Black captures a white pawn.
		}
	}

	// Update the e.p. target.
	if isPawnMove {
		fr, tr := from.Rank(), to.Rank()
		switch {
		case fr == Rank2 && tr == Rank4:
			p.ep.Set(to.Below()) // White double push.
		case fr == Rank7 && tr == Rank5:
			p.ep.Set(to.Above()) // Black double push.
		default:
			p.ep.Clear() // Other pawn move.
		}
	} else {
		p.ep.Clear() // Not a pawn move.
	}

	// If moving a king, update castling rights.
	switch held {
	case WhiteKing:
		p.cr.ClearWhite()
	case BlackKing:
		p.cr.ClearBlack()
	}

	// If moving from or to a corner square, update castling rights.
	switch {
	case from == A1 || to == A1:
		p.cr.ClearWhiteOOO()
	case from == H1 || to == H1:
		p.cr.ClearWhiteOO()
	case from == A8 || to == A8:
		p.cr.ClearBlackOOO()
	case from == H8 || to == H8:
		p.cr.ClearBlackOO()
	}

	// If promoting, update the held piece.
	if pt, ok := m.Promotion(); ok {
		held = NewPiece(p.active, pt)
	}

	// Move the piece from its initial square to its final square.
	// If castling, this changes the location of the king.
	p.board.Clear(from)
	p.board.Set(held, to)

	// If castling, move the castled rook.
	if held.Type() == King {
		switch {
		case from == E1 && to == G1: // WhiteOO
			p.board.Clear(H1)
			p.board.Set(WhiteRook, F1)
		case from == E1 && to == C1: // WhiteOOO
			p.board.Clear(A1)
			p.board.Set(WhiteRook, D1)
		case from == E8 && to == G8: // BlackOO
			p.board.Clear(H8)
			p.board.Set(BlackRook, F8)
		case from == E8 && to == C8: // BlackOOO
			p.board.Clear(A8)
			p.board.Set(BlackRook, D8)
		}
	}

	// Update the half move clock.
	if isPawnMove || isCapture {
		p.hmc = 0
	} else {
		p.hmc++
	}

	// Update the full move number.
	if p.active == Black {
		p.fmn++
	}

	// Switch sides.
	p.active = p.active.Other()
}
