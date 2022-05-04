package key

import (
	"crypto/ecdsa"
	"errors"
	"regexp"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	prkRegExp                  = regexp.MustCompile(`^[a-f0-9]{64}$`)
	addressRegExp              = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	ErrCastingPublicKeyToECDSA = errors.New("error casting public key to ECDSA")
)

// Generate returns a private key and its public address as two strings.
func Generate() (prk string, a string, err error) {
	prkECDSA, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	prk = hexutil.Encode(crypto.FromECDSA(prkECDSA))[2:]

	pbkECDSA, ok := prkECDSA.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", "", ErrCastingPublicKeyToECDSA
	}
	a = crypto.PubkeyToAddress(*pbkECDSA).Hex()
	return
}

// GenerateN returns a map of private keys and public addresses. If N <= 0, generates a single entry.
func GenerateN(N int) (map[string]string, error) {
	n := 0
	if N <= 0 {
		n = 1
	} else {
		n = N
	}

	m := map[string]string{}
	for i := 0; i < n; i++ {
		k, a, err := Generate()
		if err != nil {
			return nil, err
		}
		m[k] = a
	}
	return m, nil
}

// ValidPrivateKey tests if the key is a well formated private key
func ValidPrivateKey(k string) bool {
	return prkRegExp.Match([]byte(k))
}

// ValidPublicAddress tests if the address is a well formated public address
func ValidPublicAddress(a string) bool {
	return addressRegExp.Match([]byte(a))
}
