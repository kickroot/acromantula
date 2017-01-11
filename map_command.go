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

import "fmt"

//
// We can utilize the same logic and code for headers, settings, and params
//
type mapCommand struct {
	desc       string
	backingMap map[string]string
}

func (c *mapCommand) description() string {
	return c.desc
}

func (c *mapCommand) usage() string {
	return fmt.Sprintf("[set <key> <value>] | [unset <key>]")
}

func (c *mapCommand) exec(tokens []string, term *Term, config *configuration) {

	//
	// If only the top level command is specified, then we simply print the
	// contents
	//
	if len(tokens) == 1 {
		for _, k := range sortKeys(c.backingMap) {
			term.printf(" %v => %v\n", k, c.backingMap[k])
		}
		return
	}

	switch tokens[1] {
	case "set":
		if len(tokens) < 3 {
			term.printf("No name/value supplied, try '%s set <name> <value>'\n", tokens[0])
		} else if len(tokens) < 4 {
			term.printf("%s needs a value as well, try '%s set %s <value>'\n", tokens[2], tokens[0], tokens[2])
		} else {
			c.backingMap[tokens[2]] = tokens[3]
		}
	case "unset":
		if len(tokens) < 3 {
			term.printf("No key supplied, try '%s unset <name> [name...]'\n", tokens[0])
		} else {
			for _, key := range tokens[2:] {
				delete(c.backingMap, key)
			}
		}
	default:
		term.printf("Unknown sub-command '%s', try one of [set, unset]\n", tokens[1])
	}
}
