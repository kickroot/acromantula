package main

import (
  "strings"
  "os"
  "log"
  "golang.org/x/crypto/ssh/terminal"
)

var termState *terminal.State
var term terminal.Terminal

func restoreTerm() {
	terminal.Restore(0, termState)
}

func initTerm() {
  oldState, err := terminal.MakeRaw(0)  
  if (err != nil) {
    log.Fatal(err)
  } else {
    termState = oldState;
    term = *terminal.NewTerminal(os.Stdin, "acro >> ")
  }
}

func readTerm() ([]string, error) {
  str, err := term.ReadLine()
  if (err != nil) {
      return  nil, err
    } else {
      tokens := strings.Split(str, " ")
      trimmedTokens := make([]string, len(tokens), len(tokens))
      for index, value :=  range tokens {
        trimmedTokens[index] = strings.TrimSpace(value)
      }
      return trimmedTokens, nil;
    }
}