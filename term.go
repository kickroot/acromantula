package main

import (
  "fmt"
  "os"

  "golang.org/x/crypto/ssh/terminal"
)

func handle_keypress(b chan bool, line string, pos int, key rune) (newLine string, newPos int, ok bool) {
  switch key {
  default:
    fmt.Printf("key:[%c] pos[%d] line:[%s]\n", key, pos, line)
  case 'q':
    b <- true
  }
  return line, pos, false
}

func main() {
  done := make(chan bool)
  go func() {
    oldState, _ := terminal.MakeRaw(0)
    terminal.MakeRaw(int(os.Stdin.Fd()))
    defer terminal.Restore(0, oldState)
    n := terminal.NewTerminal(os.Stdin, ">> ")
    f := func(s string, i int, r rune) (string, int, bool) {
      return handle_keypress(done, s, i, r)
    }
    n.AutoCompleteCallback = f
    for {
      n.ReadLine()
    }
  }()
  select {
  case <-done:
    fmt.Printf("[q] - exit\n")
    close(done)
  }
}
