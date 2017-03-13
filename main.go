// Package main provides the command line utility.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cassierecher/markovian/impl"
)

var (
	inFilePath  = flag.String("in", "", `A file to train the Markov chain on. Alternatively, specify "stdin" to use standard input.`)
	outFilePath = flag.String("out", "markov.json", "A file to store the Markov chain in. File will be overwritten if it already exists.")
	order       = flag.Int("order", 2, `The order of the Markov chain.`)
)

func init() {
	flag.Parse()
}

// Implements the "help" command.
func helpCmd() {
	fmt.Fprintf(os.Stderr, `Markovian

Synopsis: markovian ARG

Args:
- help:		Display this message.
- train:	Train a Markov chain. Relevant flags: inFilePath, outFilePath, order.

Flags:
`)
	flag.PrintDefaults()
}

// Implements the "train" command.
func trainCmd() {
	// Get the data to read.
	var r io.Reader

	switch *inFilePath {
	case "":
		fmt.Fprintf(os.Stderr, "Must provide input file path.\n")
		os.Exit(1)
	case "stdin":
		r = os.Stdin
	default:
		in, err := os.Open(*inFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer in.Close()
		r = in
	}

	mc, err := impl.New(*order)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	mc.Train(r)

	// Marshal the resulting Markov chain to JSON in a file.
	out, err := os.Create(*outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer out.Close()

	e := json.NewEncoder(out)
	if err := e.Encode(mc); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding to JSON: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	// Handle args.
	args := flag.Args()
	// Validate number of args.
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Not enough args.\n")
		helpCmd()
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Too many args.\n")
		helpCmd()
		os.Exit(1)
	}

	switch args[0] {
	case "train":
		trainCmd()
	case "help":
		helpCmd()
	default:
		fmt.Printf("Unrecognized command %q.\n")
		helpCmd()
		os.Exit(1)
	}
}
