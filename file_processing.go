package processing

import (
	"bufio"
	"fmt"
	"os"
)

type TextFileProcessor struct {
	Filename string
	ErrHandler
}

type ErrHandler func(error)

func NewTextFileProcessor(file string, errHandler ErrHandler) *TextFileProcessor {
	return &TextFileProcessor{
		Filename:   file,
		ErrHandler: errHandler,
	}
}

func (p *TextFileProcessor) processFile(outChan chan interface{}) {
	file, err := os.Open(p.Filename)
	if err != nil {
		if p.ErrHandler != nil {
			p.ErrHandler(err)
		} else {
			panic(fmt.Sprintf("Error opening file: %v", err))
		}
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		outChan <- line
	}

	if err = scanner.Err(); err != nil {
		if p.ErrHandler != nil {
			p.ErrHandler(err)
		} else {
			panic(fmt.Sprintf("Error scanning file: %v", err))
		}
	}
}

func (p *TextFileProcessor) Process() <-chan interface{} {
	outChan := make(chan interface{})

	go func() {
		p.processFile(outChan)
		close(outChan)
	}()

	return outChan
}