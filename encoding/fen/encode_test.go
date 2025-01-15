package fen

import (
	"testing"

	"github.com/clfs/lento/core"
)

func TestEncode_Starting(t *testing.T) {
	p := core.NewPosition()
	want := Starting
	if got := Encode(p); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
