package key

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrCastingPublicKeyToECDSA = errors.New("error casting public key to ECDSA")

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

func GenerateN(n int) {}
