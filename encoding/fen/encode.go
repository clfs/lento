package fen

import (
	"fmt"
	"strings"

	"github.com/clfs/lento/core"
)

// Encode encodes a position using FEN.
func Encode(p core.Position) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s ", encodeBoard(p.Board()))
	fmt.Fprintf(&sb, "%c ", encodeColor(p.Active()))
	fmt.Fprintf(&sb, "%s ", encodeCastlingRights(p.CastlingRights()))
	fmt.Fprintf(&sb, "%s ", encodeEnPassantTarget(p.EnPassantTarget()))
	fmt.Fprintf(&sb, "%d ", p.HalfmoveClock())
	fmt.Fprintf(&sb, "%d", p.FullmoveNumber())

	return sb.String()
}

func encodePiece(p core.Piece) byte {
	return map[core.Piece]byte{
		core.WhitePawn:   'P',
		core.WhiteKnight: 'N',
		core.WhiteBishop: 'B',
		core.WhiteRook:   'R',
		core.WhiteQueen:  'Q',
		core.WhiteKing:   'K',
		core.BlackPawn:   'p',
		core.BlackKnight: 'n',
		core.BlackBishop: 'b',
		core.BlackRook:   'r',
		core.BlackQueen:  'q',
		core.BlackKing:   'k',
	}[p]
}

func encodeColor(c core.Color) byte {
	if c {
		return 'b'
	}
	return 'w'
}

func encodeBoard(b core.Board) string {
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

			sb.WriteByte(encodePiece(p))
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

func encodeCastlingRights(c core.CastlingRights) string {
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

func encodeEnPassantTarget(e core.EnPassantTarget) string {
	sq, ok := e.Get()
	if !ok {
		return "-"
	}

	f := 'a' + sq.File()
	r := '1' + sq.Rank()

	return fmt.Sprintf("%c%c", f, r)
}
