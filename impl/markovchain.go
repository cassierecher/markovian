// impl is the implementation of Markovian. It is isolated from the main package to help enforce encapsulation.
package impl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type frequencyGroup map[string]int

func newFrequencyGroup() map[string]int {
	return make(map[string]int)
}

// MarkovChain encapsulates a Markov chain.
// Never create a MarkovChain directly - always use the provided New function.
type MarkovChain struct {
	Order     int
	Knowledge map[string]frequencyGroup
}

// New returns a fully-initialized MarkovChain of the given order, or the first error to block initialization.
// Order must be a positive value.
func New(order int) (*MarkovChain, error) {
	if order <= 0 {
		return nil, fmt.Errorf("order must be positive (got %d)", order)
	}
	m := MarkovChain{
		Knowledge: make(map[string]frequencyGroup),
		Order:     order,
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
		if len(currRunes) != 1 && unicode.In(last, unicode.STerm) {
			// Add the word and sentence terminal rune separately.
			currWithoutLast := currRunes[:len(currRunes)-1]
			words = append(words, string(currWithoutLast), string(last))
		} else {
			words = append(words, curr)
		}

		for _, v := range words {
			// Save the information.
			m.addKnowledge(back, v)

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

// buildKey builds a map key from the given string slice.
// It joins the strings with $s as a separator, and escapes all pre-existing $s for safety.
// The map key is a string.
func buildKey(in []string) string {
	for i := range in {
		in[i] = strings.Replace(in[i], `\`, `\\`, -1)
		in[i] = strings.Replace(in[i], `$`, `\$`, -1)
	}
	return strings.Join(in, "$")
}

func (m *MarkovChain) addKnowledge(back []string, next string) {
	k := buildKey(back)
	fg, ok := m.Knowledge[k]
	if !ok {
		fg = newFrequencyGroup()
		m.Knowledge[k] = fg
	}
	fg[next]++
}
