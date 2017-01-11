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
