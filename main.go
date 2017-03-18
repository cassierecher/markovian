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
	inFilePath   = flag.String("in", "", "A file containing a Markov chain to use. Leave empty to start with a new Markov chain.")
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
- train:	Train a Markov chain. Relevant flags: in, data, out, order.

Flags:
`)
	flag.PrintDefaults()
}

// Gets the Markov chain to focus on, either from the specified file, or new if none was specified.
// Returns the Markov chain, or an error if one should occur.
func obtainMarkovChain() (*impl.MarkovChain, error) {
	// Must use address of zero value instead of nil pointer due to JSON parsing requirement.
	mc := &impl.MarkovChain{}
	if *inFilePath == "" {
		var err error
		mc, err = impl.New(*order)
		if err != nil {
			return nil, fmt.Errorf("couldn't make new Markov chain: %s", err)
		}
	} else {
		// Get data from file.
		b, err := ioutil.ReadFile(*inFilePath)
		if err != nil {
			return nil, fmt.Errorf("couldn't read input file: %s", err)
		}
		if err := json.Unmarshal(b, mc); err != nil {
			return nil, fmt.Errorf("couldn't read json: %s", err)
		}
	}
	return mc, nil
}

// Implements the "train" command.
// Returns errors, if one should occur.
func trainCmd() error {
	// Get the input Markov chain.
	mc, err := obtainMarkovChain()
	if err != nil {
		return fmt.Errorf("couldn't obtain Markov chain: %s", err)
	}

	// Get the data to read.
	var r io.Reader
	switch *dataFilePath {
	case "":
		return errors.New("must provide input file path")
	case "stdin":
		r = os.Stdin
	default:
		// Input is possibly very large; use a reader instead of ioutil convenience method.
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
