package fen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/clfs/lento/core"
)

// MustDecode is like [Decode] but panics if the FEN string is invalid.
func MustDecode(s string) core.Position {
	p, err := Decode(s)
	if err != nil {
		panic(err)
	}
	return p
}

// Decode decodes a position from FEN.
func Decode(s string) (core.Position, error) {
	fields := strings.Split(s, " ")
	if n := len(fields); n != 6 {
		return core.Position{}, fmt.Errorf("bad field count: %d", n)
	}

	board, err := DecodeBoard(fields[0])
	if err != nil {
		return core.Position{}, fmt.Errorf("bad board: %v", err)
	}

	sideToMove, err := DecodeColor(fields[1])
	if err != nil {
		return core.Position{}, err
	}

	castlingRights, err := DecodeCastlingRights(fields[2])
	if err != nil {
		return core.Position{}, err
	}

	epTarget, epExists, err := DecodeEnPassantTarget(fields[3])
	if err != nil {
		return core.Position{}, err
	}

	// TODO(clfs): Validate this more cleanly.
	if epExists {
		switch epTarget.Rank() {
		case core.Rank3:
			if sideToMove == core.White {
				return core.Position{}, fmt.Errorf("invalid e.p. square for white: %s", fields[3])
			}
		case core.Rank6:
			if sideToMove == core.Black {
				return core.Position{}, fmt.Errorf("invalid e.p. square for black: %s", fields[3])
			}
		default:
			return core.Position{}, fmt.Errorf("invalid rank for e.p. square: %s", fields[3])
		}
	}

	halfmoveClock, err := decodeHalfmoveClock(fields[4])
	if err != nil {
		return core.Position{}, err
	}

	fullmoveNumber, err := decodeFullmoveNumber(fields[5])
	if err != nil {
		return core.Position{}, err
	}

	opts := []core.PositionOption{
		core.WithBoard(board),
		core.WithSideToMove(sideToMove),
		core.WithCastlingRights(castlingRights),
		core.WithHalfmoveClock(halfmoveClock),
		core.WithFullmoveNumber(fullmoveNumber),
	}

	if epExists {
		opts = append(opts, core.WithEnPassantTarget(epTarget))
	}

	p := core.NewPosition(opts...)

	return p, nil
}

func decodePiece(s string) (core.Piece, error) {
	m := map[string]core.Piece{
		"P": core.WhitePawn,
		"N": core.WhiteKnight,
		"B": core.WhiteBishop,
		"R": core.WhiteRook,
		"Q": core.WhiteQueen,
		"K": core.WhiteKing,
		"p": core.BlackPawn,
		"n": core.BlackKnight,
		"b": core.BlackBishop,
		"r": core.BlackRook,
		"q": core.BlackQueen,
		"k": core.BlackKing,
	}

	p, ok := m[s]
	if !ok {
		return 0, fmt.Errorf("bad piece: %q", s)
	}
	return p, nil
}

// DecodeBoard decodes a board from FEN.
func DecodeBoard(s string) (core.Board, error) {
	var b core.Board

	offset := int(core.A8) // top left corner

	ranks := strings.Split(s, "/")
	if n := len(ranks); n != 8 {
		return core.Board{}, fmt.Errorf("bad rank count: %d", n)
	}

	for i, rank := range ranks {
		var numPrev bool
		for _, rn := range rank {
			switch rn {
			case '1', '2', '3', '4', '5', '6', '7', '8':
				if numPrev {
					return core.Board{}, fmt.Errorf("bad rank: %q", rank)
				}
				offset += int(rn - '0') // advance rightwards by n
				numPrev = true
			default:
				piece, err := decodePiece(string(rn))
				if err != nil {
					return core.Board{}, err
				}
				b.Set(piece, core.Square(offset))
				offset++ // advance rightwards by 1
				numPrev = false
			}
		}

		// Were all eight squares accounted for?
		if offset != 8*(8-i) {
			return core.Board{}, fmt.Errorf("bad rank: %q", rank)
		}

		offset -= 16 // advance down by 2
	}

	return b, nil
}

// DecodeColor decodes a color from FEN.
func DecodeColor(s string) (core.Color, error) {
	switch s {
	case "w":
		return core.White, nil
	case "b":
		return core.Black, nil
	default:
		return false, fmt.Errorf("bad color: %q", s)
	}
}

// DecodeCastlingRights decodes castling rights from FEN.
func DecodeCastlingRights(s string) (core.CastlingRights, error) {
	var cr core.CastlingRights

	switch s {
	case "-":
		// nothing
	case "K":
		cr.SetWhiteOO()
	case "Q":
		cr.SetWhiteOOO()
	case "k":
		cr.SetBlackOO()
	case "q":
		cr.SetBlackOOO()
	case "KQ":
		cr.SetWhiteOO()
		cr.SetWhiteOOO()
	case "Kk":
		cr.SetWhiteOO()
		cr.SetBlackOO()
	case "Kq":
		cr.SetWhiteOO()
		cr.SetBlackOOO()
	case "Qk":
		cr.SetWhiteOOO()
		cr.SetBlackOO()
	case "Qq":
		cr.SetWhiteOOO()
		cr.SetBlackOOO()
	case "kq":
		cr.SetBlackOO()
		cr.SetBlackOOO()
	case "KQk":
		cr.SetWhiteOO()
		cr.SetWhiteOOO()
		cr.SetBlackOO()
	case "KQq":
		cr.SetWhiteOO()
		cr.SetWhiteOOO()
		cr.SetBlackOOO()
	case "Kkq":
		cr.SetWhiteOO()
		cr.SetBlackOO()
		cr.SetBlackOOO()
	case "Qkq":
		cr.SetWhiteOOO()
		cr.SetBlackOO()
		cr.SetBlackOOO()
	case "KQkq":
		cr.SetWhiteOO()
		cr.SetWhiteOOO()
		cr.SetBlackOO()
		cr.SetBlackOOO()
	default:
		return core.CastlingRights{}, fmt.Errorf("invalid castling rights: %q", s)
	}

	return cr, nil
}

// DecodeEnPassantTarget decodes an en passant target from FEN.
func DecodeEnPassantTarget(s string) (core.Square, bool, error) {
	// TODO(clfs): Write this more cleanly.
	if s == "-" {
		return 0, false, nil
	}

	if len(s) != 2 {
		return 0, false, fmt.Errorf("invalid e.p. target: %q", s)
	}

	f := core.File(s[0] - 'a')
	r := core.Rank(s[1] - '1')

	if f > core.FileH || r > core.Rank8 {
		return 0, false, fmt.Errorf("invalid e.p. target: %q", s)
	}

	sq := core.NewSquare(f, r)

	return sq, true, nil
}

var halfmoveClockRegexp = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)

func decodeHalfmoveClock(s string) (int, error) {
	if !halfmoveClockRegexp.MatchString(s) {
		return 0, fmt.Errorf("bad halfmove clock %q", s)
	}
	return strconv.Atoi(s)
}

var fullmoveNumberRegexp = regexp.MustCompile(`^([1-9][0-9]*)$`)

func decodeFullmoveNumber(s string) (int, error) {
	if !fullmoveNumberRegexp.MatchString(s) {
		return 0, fmt.Errorf("bad fullmove number %q", s)
	}
	return strconv.Atoi(s)
}
