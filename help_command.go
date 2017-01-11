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

type helpCommand struct{}

func (c *helpCommand) exec(tokens []string, term *Term, config *configuration) {
	if len(tokens) == 1 {
		term.printf("For help with a command, try 'help <cmd>', where <cmd> is one of:\n")

		for _, key := range sortCommands(commands) {
			term.printf("  %s - %s\n", key, commands[key].description())
		}
	} else if len(tokens) == 2 {
		cmd := commands[tokens[1]]
		if cmd == nil {
			term.printf("Unknown command: %s\n", tokens[1])
		} else {
			term.printf("%s: %s\n", tokens[1], cmd.description())
			term.printf("Usage: %s %s\n", tokens[1], cmd.usage())
		}
	} else {
		term.printf("%s\n", c.usage())
	}
}

func (c *helpCommand) description() string {
	return "Displays basic usage information on all commands."
}

func (c *helpCommand) usage() string {
	return "help <command>"
}
