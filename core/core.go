// Package core implements basic chess functionality.
package core

// A Color represents either [White] or [Black].
type Color bool

// [Color] constants.
const (
	White Color = false
	Black Color = true
)

// Other returns the opposite color.
func (c Color) Other() Color {
	return !c
}

// Uint64 returns 0 for [White] and 1 for [Black].
func (c Color) Uint64() uint64 {
	if c {
		return 1
	}
	return 0
}

// A PieceType represents a type of piece.
type PieceType uint8

// [PieceType] constants.
const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// A Piece represents a chess piece.
type Piece uint8

// [Piece] constants.
const (
	WhitePawn Piece = iota
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

// NewPiece returns a new piece.
func NewPiece(c Color, pt PieceType) Piece {
	if c {
		return Piece(pt) + 6
	}
	return Piece(pt)
}

// Color returns the piece's color.
func (p Piece) Color() Color {
	return p >= 6
}

// Type returns the piece's type.
func (p Piece) Type() PieceType {
	return PieceType(p % 6)
}

// A File is a column on the board.
type File uint8

// [File] constants.
const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

// A Rank is a row on the board.
type Rank uint8

// [Rank] constants. Note that [Rank1] is equal to 0.
const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

// Above returns the rank above r.
// It is invalid to call Above if r is [Rank8].
func (r Rank) Above() Rank {
	return r + 1
}

// Below returns the rank below r.
// It is invalid to call Below if r is [Rank1].
func (r Rank) Below() Rank {
	return r - 1
}

// A Square is a location on the board.
type Square uint16

// [Square] constants.
const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// NewSquare returns a new square.
func NewSquare(f File, r Rank) Square {
	return Square(f) + Square(r)*8
}

// File returns the file that the square lies on.
func (s Square) File() File {
	return File(s % 8)
}

// Rank returns the rank that the square lies on.
func (s Square) Rank() Rank {
	return Rank(s / 8)
}

// Above returns the square above s.
// It is invalid to call Above if s is on the eighth rank.
func (s Square) Above() Square {
	return s + 8
}

// Below returns the square below s.
// It is invalid to call Below if s is on the first rank.
func (s Square) Below() Square {
	return s - 8
}
