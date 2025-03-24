// Package fen implements encoding and decoding Forsyth-Edwards Notation.
//
// This package is compliant with "Standard: Portable Game Notation
// Specification and Implementation Guide", revision 1994.03.12, §16.1.3.4 and
// will encode an en passant target square "if and only if the last move was a
// pawn advance of two squares [...] even if there is no pawn of the opposing
// side that may immediately execute the en passant capture".
//
// TODO(clfs): Add Go examples.
//
// TODO(clfs): Document which situations cause decoding errors and which don't.
package fen

// Starting is the FEN string for the starting position.
const Starting = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
