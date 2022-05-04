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
