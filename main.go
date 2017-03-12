// Package main provides the command line utility.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cassierecher/markovian/impl"
)

var inFilePath = flag.String("in", "", `A file to train the Markov chain on. Alternatively, specify "stdin" to use standard input.`)
var order = flag.Int("order", 2, `The order of the Markov chain.`)

func init() {
	flag.Parse()
}

func trainCmd() {
	// Get the data to read.
	var r io.Reader

	switch *inFilePath {
	case "":
		fmt.Println("Must provide input file path.")
		os.Exit(1)
	case "stdin":
		r = os.Stdin
	default:
		in, err := os.Open(*inFilePath)
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		defer in.Close()
		r = in
	}

	mc := impl.New()
	mc.Train(r, *order)
	fmt.Printf("%+v\n", mc)
}

func main() {
	// Handle args.
	if len(flag.Args()) > 1 {
		fmt.Printf("Too many arguments.")
		os.Exit(1)
	}

	trainCmd()

}
