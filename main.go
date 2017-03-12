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

// Implements the "help" command.
func helpCmd() {
	fmt.Println(`Markovian

Synopsis: markovian ARG

Args:
- help:		Display this message.
- train:	Train a Markov chain. Relevant flags: inFilePath, order.
`)
}

// Implements the "train" command.
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
	args := flag.Args()
	// Validate number of args.
	if len(args) < 1 {
		fmt.Println("Not enough args.")
		helpCmd()
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Too many args.")
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
