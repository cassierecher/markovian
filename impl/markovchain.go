// impl is the implementation of Markovian. It is isolated from the main package to help enforce encapsulation.
package impl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// MarkovChain encapsulates a Markov chain.
// Never create a MarkovChain directly - always use the provided New function.
type MarkovChain struct {
	Lessons []lesson
	Order   int
}

// lesson internally represents a single iteration of training - a set of words, and the word to follow.
type lesson struct {
	Back []string
	Next string
}

// New returns a fully-initialized MarkovChain of the given order, or the first error to block initialization.
// Order must be a positive value.
func New(order int) (*MarkovChain, error) {
	if order <= 0 {
		return nil, fmt.Errorf("order must be positive (got %d)", order)
	}
	m := MarkovChain{
		Order: order,
	}
	return &m, nil
}

// Train takes a reader and trains the Markov chain with the data in the reader.
// It returns an error, should one occur.
func (m *MarkovChain) Train(r io.Reader) error {
	// Validate input.
	if r == nil {
		return errors.New("got r = nil, want non-nil")
	}

	// Construct a word scanner.
	wordScnr := bufio.NewScanner(r)
	wordScnr.Split(bufio.ScanWords)

	back := make([]string, m.Order, m.Order)

	// Scan until the scanner won't scan anymore.
	for wordScnr.Scan() {
		// Interpret the scan value as a string.
		curr := string(wordScnr.Bytes())

		var words []string

		// Check if the last rune is a sentence terminator, and that it's not just a sentence terminator.
		currRunes := []rune(curr)
		last := currRunes[len(currRunes)-1]
		if len(currRunes) != 1 && (last == '?' || last == '!' || last == '.') { // TODO: handle Unicode sentence terminators more generally
			// Add the word and sentence terminal rune separately.
			currWithoutLast := currRunes[:len(currRunes)-1]
			words = append(words, string(currWithoutLast), string(last))
		} else {
			words = append(words, curr)
		}

		for _, v := range words {
			// Save the information in a lesson.
			m.Lessons = append(m.Lessons, lesson{
				Back: back,
				Next: v,
			})

			// Update state for the next iteration.
			back = append(back[1:], v)
		}
	}
	// EOF is not recognized as an error here, so we can just check for the presence of any error.
	if err := wordScnr.Err(); err != nil {
		return fmt.Errorf("encountered problem while scanning words: %s", err)
	}

	return nil
}
