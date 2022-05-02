package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, '\t', 0)

	prk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	prkBytes := crypto.FromECDSA(prk)
	fmt.Fprintf(w, "🔑 private key:\t\u001b[1;32m%s\u001b[0m\n", hexutil.Encode(prkBytes)[2:])

	pbk := prk.Public()
	pbkECDSA, ok := pbk.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	a := crypto.PubkeyToAddress(*pbkECDSA).Hex()
	fmt.Fprintf(w, "🚩 public address:\t\u001b[44m%s\u001b[0m\n", a)
	w.Flush()
}
