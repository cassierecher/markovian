// Package main provides the command line utility.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cassierecher/markovian/impl"
)

var (
	dataFilePath = flag.String("data", "", `A file of data to train the Markov chain on. Alternatively, specify "stdin" to use standard input.`)
	outFilePath  = flag.String("out", "markov.json", "A file to store the Markov chain in. File will be overwritten if it already exists.")
	order        = flag.Int("order", 2, `The order of the Markov chain.`)
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
- train:	Train a Markov chain. Relevant flags: data, out, order.

Flags:
`)
	flag.PrintDefaults()
}

// Implements the "train" command.
// Returns errors, if one should occur.
func trainCmd() error {
	// Get the input Markov chain.
	mc, err := impl.New(*order)
	if err != nil {
		return fmt.Errorf("couldn't get input Markov chain: %s", err)
	}

	// Get the data to read.
	var r io.Reader
	switch *dataFilePath {
	case "":
		return errors.New("must provide input file path")
	case "stdin":
		r = os.Stdin
	default:
		in, err := os.Open(*dataFilePath)
		if err != nil {
			return fmt.Errorf("couldn't open file: %s", err)
		}
		defer in.Close()
		r = in
	}

	// Perform training.
	mc.Train(r)

	// Marshal the resulting Markov chain to JSON in a file.
	b, err := json.Marshal(mc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*outFilePath, b, 0600)
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
		if err := trainCmd(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case "help":
		helpCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized command %q.\n", args[0])
		helpCmd()
		os.Exit(1)
	}
}
