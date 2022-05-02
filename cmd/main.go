package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	key "github.com/fairhive-labs/ethkeygen/pkg"
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, '\t', 0)
	prk, a, err := key.Generate()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "🔑 private key:\t\u001b[1;34m%s\u001b[0m\n", prk)
	fmt.Fprintf(w, "🚩 public address:\t\u001b[42m%s\u001b[0m\n", a)
	w.Flush()
}
