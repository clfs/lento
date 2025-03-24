package fen

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readFENFile(t testing.TB, name string) []string {
	t.Helper()

	var res []string

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Ignore blank lines and comments.
		}

		res = append(res, line)
	}

	if err := s.Err(); err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("empty file")
	}

	return res
}

func TestDecode_Valid(t *testing.T) {
	for _, s := range readFENFile(t, "testdata/valid.fen") {
		if _, err := Decode(s); err != nil {
			t.Errorf("%q: error: %v", s, err)
		}
	}
}

func TestDecode_Invalid(t *testing.T) {
	for _, s := range readFENFile(t, "testdata/invalid.fen") {
		if _, err := Decode(s); err == nil {
			t.Errorf("%q: no error", s)
		}
	}
}

func FuzzRoundTrip(f *testing.F) {
	corpuses, err := filepath.Glob("testdata/*.fen")
	if err != nil {
		f.Fatal(err)
	}

	for _, c := range corpuses {
		for _, s := range readFENFile(f, c) {
			f.Add(s)
		}
	}

	f.Fuzz(func(t *testing.T, s string) {
		p, err := Decode(s)
		if err != nil {
			return
		}

		s2 := Encode(p)
		if s != s2 {
			t.Fatalf("changed in round trip: %q -> %q", s, s2)
		}
	})
}
