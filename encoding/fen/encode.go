package fen

import (
	"fmt"
	"strings"

	"github.com/clfs/lento/core"
)

// Encode encodes a position into FEN.
func Encode(p core.Position) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s ", EncodeBoard(p.Board()))
	fmt.Fprintf(&b, "%s ", EncodeColor(p.SideToMove()))
	fmt.Fprintf(&b, "%s ", EncodeCastlingRights(p.CastlingRights()))
	fmt.Fprintf(&b, "%s ", EncodeEnPassantTarget(p.EnPassantTarget()))
	fmt.Fprintf(&b, "%d ", p.HalfmoveClock())
	fmt.Fprintf(&b, "%d", p.FullmoveNumber())

	return b.String()
}

// EncodePiece encodes a piece into FEN.
func EncodePiece(p core.Piece) string {
	m := map[core.Piece]string{
		core.WhitePawn:   "P",
		core.WhiteKnight: "N",
		core.WhiteBishop: "B",
		core.WhiteRook:   "R",
		core.WhiteQueen:  "Q",
		core.WhiteKing:   "K",
		core.BlackPawn:   "p",
		core.BlackKnight: "n",
		core.BlackBishop: "b",
		core.BlackRook:   "r",
		core.BlackQueen:  "q",
		core.BlackKing:   "k",
	}
	return m[p]
}

// EncodeColor encodes a color into FEN.
func EncodeColor(c core.Color) string {
	if c == core.White {
		return "w"
	}
	return "b"
}

// EncodeBoard encodes a board into FEN.
func EncodeBoard(b core.Board) string {
	var sb strings.Builder

	for r := core.Rank8; r <= core.Rank8; r-- {
		gap := 0
		for f := core.FileA; f <= core.FileH; f++ {
			p, ok := b.Get(core.NewSquare(f, r))

			// Empty square?
			if !ok {
				gap++
				continue
			}

			// End of gap?
			if gap > 0 {
				fmt.Fprintf(&sb, "%d", gap)
				gap = 0
			}

			sb.WriteString(EncodePiece(p))
		}

		// Row ends in gap?
		if gap > 0 {
			fmt.Fprintf(&sb, "%d", gap)
		}

		// Row divider needed?
		if r != core.Rank1 {
			sb.WriteByte('/')
		}
	}

	return sb.String()
}

// EncodeCastlingRights encodes castling rights into FEN.
func EncodeCastlingRights(c core.CastlingRights) string {
	var sb strings.Builder

	if c.GetWhiteOO() {
		sb.WriteByte('K')
	}
	if c.GetWhiteOOO() {
		sb.WriteByte('Q')
	}
	if c.GetBlackOO() {
		sb.WriteByte('k')
	}
	if c.GetBlackOOO() {
		sb.WriteByte('q')
	}

	if sb.Len() == 0 {
		return "-"
	}
	return sb.String()
}

// EncodeEnPassantTarget encodes an en passant target into FEN.
func EncodeEnPassantTarget(e core.EnPassantTarget) string {
	sq, ok := e.Get()
	if !ok {
		return "-"
	}

	f := 'a' + sq.File()
	r := '1' + sq.Rank()

	return fmt.Sprintf("%c%c", f, r)
}
