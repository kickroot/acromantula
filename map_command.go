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
