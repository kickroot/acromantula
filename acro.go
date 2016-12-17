package main

import (
	"fmt"
	"strings"
)

var root string
var term *Term
var headers *Headers

func main() {
	term = createTerm(0)
	headers = createHeaders()

	fmt.Printf("Hit Ctrl+D to quit\r\n")

	defer term.restoreTerm()

	setRoot("https://api.sourceclear.com")

	for {
		tokens, err := term.readline()

		if err != nil {
			fmt.Printf("\r\nExiting....%v\r\n", err)
			break
		}

		if len(tokens) == 0 || len(tokens[0]) == 0 {
			continue
		}

		switch tokens[0] {
		case "headers":
			handleHeaders(tokens)
		case "root":
			handleRoot(tokens)
		case "get":
			handleGet(tokens)
		default:
			fmt.Printf("Unknown command, %v\r\n", tokens[0])
		}
	}
}

func setRoot(str string) {
	root = str
}

func handleRoot(tokens []string) {

	if len(tokens) == 1 {
		fmt.Printf("%v\r\n", root)
	} else {
		setRoot(strings.TrimSpace(tokens[1]))
	}
}

func handleGet(tokens []string) {

	url := root
	if len(tokens) > 1 {
		url = root + tokens[1]
	}

	performGet(term, url, headers)
}

func handleHeaders(tokens []string) {

	//
	// In the case of just 'headers'
	//
	if len(tokens) == 1 {
		fmt.Printf("Headers\r\n")
		for k, v := range headers.all() {
			fmt.Printf("%v => %v\r\n", k, v)
		}
		return
	}

	switch tokens[1] {
	case "set":
		if len(tokens) < 4 {
			fmt.Printf("%v is missing a value\r\n")
		} else {
			headers.add(tokens[2], tokens[3])
		}
	default:
		fmt.Printf("Unknown option %v\r\n", tokens[1])
	}

}
