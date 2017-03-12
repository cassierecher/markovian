package impl

import (
	"bufio"
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
