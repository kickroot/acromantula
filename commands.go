package main

import "fmt"

//
// Command encapsulate the behavior of a command, It is expected that
// commands may contain sub-commands.
//
type command interface {

	//
	// exec, as the name implies, executes the command, given the
	// supplied context.  tokens must always contain at least one
	// item (the top level command).
	//
	// Tokens will always be the comoplete slice of parsed tokens
	exec(tokens []string, term *Term, config *configuration)

	usage() string

	description() string
}
