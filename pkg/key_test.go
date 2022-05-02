package key

import (
	"fmt"
	"regexp"
	"testing"
)

var (
	prkRegExp     = regexp.MustCompile(`^[a-f0-9]{64}$`)
	addressRegExp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
)

func TestGenerate(t *testing.T) {
	k, a, err := Generate()
	if err != nil {
		t.Errorf("error generating private key / public address: %v\n", err)
		t.FailNow()
	}

	if !prkRegExp.Match([]byte(k)) {
		t.Errorf("%q is not a private key !\n", k)
		t.FailNow()
	}

	if !addressRegExp.Match([]byte(a)) {
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
				if !prkRegExp.Match([]byte(k)) {
					t.Errorf("%q is not a private key !\n", k)
					t.FailNow()
				}

				if !addressRegExp.Match([]byte(a)) {
					t.Errorf("%q is not a public address !\n", a)
					t.FailNow()
				}
			}
		})
	}
}
