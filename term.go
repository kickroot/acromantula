package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Term is the abstraction for terminal I/O
type Term struct {
	termState *terminal.State
	term      terminal.Terminal
	fd        int
}

func createTerm(fd int) *Term {
	t := new(Term)
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	} else {
		t.fd = fd
		t.termState = oldState
		t.term = *terminal.NewTerminal(os.Stdin, "acro >> ")
		return t
	}
	return nil
}

func (t *Term) restoreTerm() {
	terminal.Restore(t.fd, t.termState)
}

func (t *Term) writeString(str string) {
	t.term.Write([]byte(str))
}

func (t *Term) writeBytes(bytes []byte) {
	t.term.Write(bytes)
}

func (t *Term) readline() ([]string, error) {
	str, err := t.term.ReadLine()
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(str, " ")
	trimmedTokens := make([]string, len(tokens), len(tokens))
	for index, value := range tokens {
		trimmedTokens[index] = strings.TrimSpace(value)
	}
	return trimmedTokens, nil
}
