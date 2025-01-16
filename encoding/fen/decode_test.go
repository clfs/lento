package fen

import (
	"testing"
)

func TestRoundTrip_Starting(t *testing.T) {
	p, err := Decode(Starting)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if got := Encode(p); got != Starting {
		t.Errorf("want %q, got %q", Starting, got)
	}
}

func TestRoundTrip_e2e4(t *testing.T) {
	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	p, err := Decode(want)
	if err != nil {
		t.Errorf("decode error: %v", err)
	}

	if got := Encode(p); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
