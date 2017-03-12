package impl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type MarkovChain struct {
	lessons []lesson
}

type lesson struct {
	back []string
	next string
}

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

	for wordScnr.Scan() {
		curr := string(wordScnr.Bytes())

		// Save the information in a lesson.
		m.lessons = append(m.lessons, lesson{
			back: back,
			next: curr,
		})

		// Update state for the next iteration.
		back = append(back[1:], curr)
	}
	if err := wordScnr.Err(); err != nil {
		return fmt.Errorf("encountered problem while scanning words: %s", err)
	}

	return nil
}
