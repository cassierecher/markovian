package impl

import (
	"bufio"
	"fmt"
	"io"
)

type MarkovChain struct {
	lessons []Lesson
}

type Lesson struct {
	Back []string
	Next string
}

func (m *MarkovChain) Train(r io.Reader, order int) error {
	// Construct a word scanner.
	wordScnr := bufio.NewScanner(r)
	wordScnr.Split(bufio.ScanWords)

	back := make([]string, order, order)

	for wordScnr.Scan() {
		curr := string(wordScnr.Bytes())

		// Save the information in a lesson.
		m.lessons = append(m.lessons, Lesson{
			Back: back,
			Next: curr,
		})

		// Update state for the next iteration.
		back = append(back[1:], curr)
	}
	if err := wordScnr.Err(); err != nil {
		return fmt.Errorf("encountered problem while scanning words: %s", err)
	}

	return nil
}
