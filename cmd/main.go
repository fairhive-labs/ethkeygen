package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	key "github.com/fairhive-labs/ethkeygen/pkg"
)

func main() {
	n := flag.Int("n", 1, "set the number of private key / public address that will be generated")
	flag.Parse()
	*n = min(*n, 100) // max 100
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, '\t', 0)
	fmt.Fprintf(w, "\n\u001b[42m** Generating %d private key(s) / public address(es) **\u001b[0m\n", *n)
	for i := range *n {
		prk, a, err := key.Generate()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "\n%d\t🔑 private key:\t\u001b[44m%s\u001b[0m\n", i+1, prk)
		fmt.Fprintf(w, "\t🚩 public address:\t\u001b[1;34m%s\u001b[0m\n", a)
	}
	w.Flush()
}
