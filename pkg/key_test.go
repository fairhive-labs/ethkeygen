package key

import (
	"fmt"
	"testing"
)

func TestValidPrivateKey(t *testing.T) {
	tt := []struct {
		k string
		v bool
	}{
		{"24c2075a0d2b0ead774eb85248f0144f785344a46359b7b0298ad7fd1c62cc0c", true},
		{"24c2075a0d2b0ead77", false},
		{"not-a-key", false},
	}
	for _, tc := range tt {
		t.Run(tc.k, func(t *testing.T) {
			if ok := ValidPrivateKey(tc.k); ok != tc.v {
				t.Errorf("incorrect private key validation, got %v, want %v\n", ok, tc.v)
				t.FailNow()
			}
		})
	}
}

func TestValidPublicAddress(t *testing.T) {
	tt := []struct {
		a string
		v bool
	}{
		{"0xa704ae62578d84A71958E33baa86F76d853C6E37", true},
		{"0xa704ae6257853C6E37", false},
		{"not-an-address", false},
	}
	for _, tc := range tt {
		t.Run(tc.a, func(t *testing.T) {
			if ok := ValidPublicAddress(tc.a); ok != tc.v {
				t.Errorf("incorrect public address validation, got %v, want %v\n", ok, tc.v)
				t.FailNow()
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	k, a, err := Generate()
	if err != nil {
		t.Errorf("error generating private key / public address: %v\n", err)
		t.FailNow()
	}

	if !ValidPrivateKey(k) {
		t.Errorf("%q is not a private key !\n", k)
		t.FailNow()
	}

	if !ValidPublicAddress(a) {
		t.Errorf("%q is not a public address !\n", a)
		t.FailNow()
	}
}

func TestGenerateN(t *testing.T) {
	tt := []struct {
		v, exp int
	}{
		{-1, 1},
		{2, 2},
		{5, 5},
		{10, 10},
		{100, 100},
		{1000, 1000},
		{10000, 10000},
		{0, 1},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc.v), func(t *testing.T) {
			m, err := GenerateN(tc.v)
			if err != nil {
				t.Errorf("error generating %d private key / public address: %v\n", tc.v, err)
				t.FailNow()
			}

			if len(m) != tc.exp {
				t.Errorf("incorrect map length, got %d, want %d\n", len(m), tc.exp)
				t.FailNow()
			}
			for k, a := range m {
				if !ValidPrivateKey(k) {
					t.Errorf("%q is not a private key !\n", k)
					t.FailNow()
				}

				if !ValidPublicAddress(a) {
					t.Errorf("%q is not a public address !\n", a)
					t.FailNow()
				}
			}
		})
	}
}

func TestSignMessage(t *testing.T) {
	tt := []struct {
		name string
		k, m string
		s    string
		err  error
	}{
		{
			"valid_1",
			"fe94806c6880c4271825152ed2ac0defa04d61b478892c4fbefae2c575a49612",
			"awesome super message to sign",
			"0xc97c389e5f120b1b4b189159b31f45c353a403c0f50a624bb2f44006a28cdaa03b23eb3d92bad05bd63a745c6f5a394c9562cf54cbbb333e190e414c7d98e9bc01",
			nil},
		{
			"valid_2",
			"fe94806c6880c4271825152ed2ac0defa04d61b478892c4fbefae2c575a49612",
			"another super message to sign",
			"0xf3fb762b0dccfc57a030741bbb998eae22e40b8c5729fa372fd27599594b71193f08598a8982adf95b12bf4a8206a4f9e79d16b24747c750fcb9c7fd675d932c01",
			nil},
		{
			"empty message",
			"fe94806c6880c4271825152ed2ac0defa04d61b478892c4fbefae2c575a49612",
			"",
			"0x324ccd8627137f5f8782280fc635e7d01fa02f6353aab68a9eeb9f7be999a8236ec48de7650ec6bd44f9605cac560c0cf79a8e06c8c4e460c9a816a680a728f300",
			nil},
		{
			"empty key",
			"",
			"another super message to sign",
			"0xf3fb762b0dccfc57a030741bbb998eae22e40b8c5729fa372fd27599594b71193f08598a8982adf95b12bf4a8206a4f9e79d16b24747c750fcb9c7fd675d932c01",
			ErrConvertingPrivateKeyToECDSA},
		{
			"invalid key",
			"123456",
			"another super message to sign",
			"0xf3fb762b0dccfc57a030741bbb998eae22e40b8c5729fa372fd27599594b71193f08598a8982adf95b12bf4a8206a4f9e79d16b24747c750fcb9c7fd675d932c01",
			ErrConvertingPrivateKeyToECDSA},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := SignMessage(tc.k, tc.m)
			if err != tc.err {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.FailNow()
			}
			if err == nil && s != tc.s {
				t.Errorf("incorrect signature, got %v, want %v", s, tc.s)
			}
		})
	}
}
