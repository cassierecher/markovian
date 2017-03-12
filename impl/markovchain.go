// impl is the implementation of Markovian. It is isolated from the main package to help enforce encapsulation.
package impl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// MarkovChain encapsulates a Markov chain.
type MarkovChain struct {
	lessons []lesson
}

// lesson internally represents a single iteration of training - a set of words, and the word to follow.
type lesson struct {
	back []string
	next string
}

// New returns a fully-initialized MarkovChain.
func New() *MarkovChain {
	return new(MarkovChain)
}

// Train takes a reader, and an order, and trains the Markov chain to the given order with the data in the reader.
// Order must be a positive value.
// It returns an error, should one occur.
func (m *MarkovChain) Train(r io.Reader, order int) error {
	// Validate input.
	if r == nil {
		return errors.New("got r = nil, want non-nil")
	}
	if order <= 0 {
		return fmt.Errorf("order must be positive (got %d)", order)
	}

	// Construct a word scanner.
	wordScnr := bufio.NewScanner(r)
	wordScnr.Split(bufio.ScanWords)

	back := make([]string, order, order)

	// Scan until the scanner won't scan anymore.
	for wordScnr.Scan() {
		// Interpret the scan value as a string.
		curr := string(wordScnr.Bytes())

		// Save the information in a lesson.
		m.lessons = append(m.lessons, lesson{
			back: back,
			next: curr,
		})

		// Update state for the next iteration.
		back = append(back[1:], curr)
	}
	// EOF is not recognized as an error here, so we can just check for the presence of any error.
	if err := wordScnr.Err(); err != nil {
		return fmt.Errorf("encountered problem while scanning words: %s", err)
	}

	return nil
}
