package core

type positionOptions struct {
	board      Board
	sideToMove Color
	cr         CastlingRights
	ep         enPassantTarget
	hmc        int
	fmn        int
	moves      []Move
}

// PositionOption configures the creation of a new position.
//
// TODO(clfs): Reword.
type PositionOption interface {
	apply(*positionOptions)
}

type boardOption Board

func (b boardOption) apply(opts *positionOptions) {
	opts.board = Board(b)
}

func WithBoard(b Board) PositionOption {
	return boardOption(b)
}

type sideToMoveOption Color

func (s sideToMoveOption) apply(opts *positionOptions) {
	opts.sideToMove = Color(s)
}

func WithSideToMove(c Color) PositionOption {
	return sideToMoveOption(c)
}

type castlingRightsOption CastlingRights

func (c castlingRightsOption) apply(opts *positionOptions) {
	opts.cr = CastlingRights(c)
}

func WithCastlingRights(cr CastlingRights) PositionOption {
	return castlingRightsOption(cr)
}

type enPassantOption enPassantTarget

func (e enPassantOption) apply(opts *positionOptions) {
	opts.ep = enPassantTarget(e)
}

func WithEnPassantTarget(s Square) PositionOption {
	var ep enPassantTarget
	ep.Set(s)
	return enPassantOption(ep)
}

type halfmoveClockOption int

func (h halfmoveClockOption) apply(opts *positionOptions) {
	opts.hmc = int(h)
}

func WithHalfmoveClock(n int) PositionOption {
	return halfmoveClockOption(n)
}

type fullmoveNumberOption int

func (f fullmoveNumberOption) apply(opts *positionOptions) {
	opts.fmn = int(f)
}

func WithFullmoveNumber(n int) PositionOption {
	return fullmoveNumberOption(n)
}

type movesOption []Move

func (m movesOption) apply(opts *positionOptions) {
	opts.moves = []Move(m)
}

func WithMoves(m []Move) PositionOption {
	return movesOption(m)
}
