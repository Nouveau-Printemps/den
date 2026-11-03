package auth

import "testing"

func FuzzArgon2id(f *testing.F) {
	f.Add("hello")
	fuzz(f, Argon2id{
		T: 1,
		M: 8 * 1024,
		P: 1,
	})
}

func FuzzScrypt(f *testing.F) {
	f.Add("hello")
	fuzz(f, Scrypt{
		N: 2 << 10,
		R: 8,
		P: 1,
	})
}

func FuzzBcrypt(f *testing.F) {
	f.Add("hello")
	fuzz(f, Bcrypt{
		Cost: 10,
	})
}

func fuzz(f *testing.F, algo Password) {
	_, ok := algo.(Bcrypt)
	var mp func(string) string
	if ok {
		mp = func(s string) string {
			if len(s) <= 72 {
				return s
			} else {
				return s[:72]
			}
		}
	}
	f.Fuzz(func(t *testing.T, s string) {
		if mp != nil {
			s = mp(s)
		}
		hash, err := algo.Hash(s)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(hash)
		ok, err := algo.Compare(hash, s)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("invalid hash for", s)
		}
	})
}
