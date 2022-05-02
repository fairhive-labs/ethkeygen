package main

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	prk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	prkBytes := crypto.FromECDSA(prk)
	fmt.Printf("private key: %s\n", hexutil.Encode(prkBytes)[2:])

	pbk := prk.Public()
	pbkECDSA, ok := pbk.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	a := crypto.PubkeyToAddress(*pbkECDSA).Hex()
	fmt.Printf("public address: %s\n", a)
}
