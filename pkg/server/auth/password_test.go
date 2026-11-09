package auth

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzArgon2id(f *testing.F) {
	fuzz(f, Argon2id{
		T: 1,
		M: 8 * 1024,
		P: 1,
	})
}

func FuzzScrypt(f *testing.F) {
	fuzz(f, Scrypt{
		N: 2 << 10,
		R: 8,
		P: 1,
	})
}

func FuzzBcrypt(f *testing.F) {
	fuzz(f, Bcrypt{
		Cost: 10,
	})
}

func fuzz(f *testing.F, algo Password) {
	f.Add("hello", "world")
	var mp func(string) string
	switch algo.(type) {
	case Bcrypt:
		mp = func(s string) string {
			if len(s) <= 72 {
				return s
			} else {
				return s[:72]
			}
		}
	case Scrypt:
		mp = func(s string) string {
			return strings.ReplaceAll(s, "\x00", " ")
		}
	}
	f.Fuzz(func(t *testing.T, s, s2 string) {
		if !utf8.ValidString(s) {
			return
		}
		if !utf8.ValidString(s2) {
			return
		}
		if mp != nil {
			s = mp(s)
			s2 = mp(s2)
		}
		hash, err := algo.Hash(s)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("pass1 %q", s)
		t.Logf("pass2 %q", s2)
		t.Log(hash)
		ok, err := algo.Compare(hash, s)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("invalid hash for", s)
		}
		ok, err = algo.Compare(hash, s2)
		if err != nil {
			t.Fatal(err)
		}
		if ok != (s == s2) {
			t.Error("invalid check")
		}
	})
}
