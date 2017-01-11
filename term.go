/*
Copyright 2017 Jason Nichols (jason@kickroot.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	keyEscape = 27
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

func (t *Term) setPrompt(prompt string) {
	t.term.SetPrompt(fmt.Sprintf("%v >> ", prompt))
}

func (t *Term) printf(str string, args ...interface{}) {
	t.term.Write([]byte(fmt.Sprintf(str, args...)))
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
	return t.tokenize(str), nil
}

func (t *Term) bright() {
	t.term.Write([]byte{keyEscape, '[', '0', '1', 'm'})
}

func (t *Term) dim() {
	t.term.Write([]byte{keyEscape, '[', '0', '2', 'm'})
}

func (t *Term) underscore() {
	t.term.Write([]byte{keyEscape, '[', '0', '4', 'm'})
}

func (t *Term) reset() {
	t.term.Write([]byte{keyEscape, '[', '0', '0', 'm'})
}

func (t *Term) tokenize(str string) []string {

	// Final tokenized set of strings.  5 is a pretty middle of the road choice
	tokens := make([]string, 0, 5)

	// Used to build the intermediate token
	buffer := bytes.NewBuffer(make([]byte, 0, 0))

	isDoubleQuoted := false
	isEscaped := false

	for _, rune := range str {
		//
		// If we are in escaped mode, write the previous character
		// literally.
		//
		if isEscaped {
			buffer.WriteRune(rune)
			isEscaped = false
			continue
		}

		char := fmt.Sprintf("%c", rune)

		if char == "\\" {
			isEscaped = true
			continue
		}

		if char == "\"" {
			isDoubleQuoted = !isDoubleQuoted
			continue
		}

		//
		// We only care if we are not in double quotes
		//
		if unicode.IsSpace(rune) && !isDoubleQuoted {
			tokens = append(tokens, buffer.String())
			buffer.Reset()
		} else {
			buffer.WriteRune(rune)
		}

	}

	//
	// At this point we should certainly be out of any quoted context
	//
	if isDoubleQuoted {
		t.writeString("Error, double quotes don't seem to match up\n")
		return []string{}
	}

	//
	// Push the remaining string
	//
	tokens = append(tokens, buffer.String())

	return tokens
}
