package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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
}

//
// We can utilize the same logic and code for headers, settings, and params
//
type mapCommand struct {
	backingMap map[string]string
}

func (c *mapCommand) exec(tokens []string, term *Term, config *configuration) {

	//
	// If only the top level command is specified, then we simply print the
	// contents
	//
	if len(tokens) == 1 {
		for k, v := range c.backingMap {
			term.printf(" %v => %v\n", k, v)
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

type httpCommand struct {
	method string
}

func (c *httpCommand) exec(tokens []string, term *Term, config *configuration) {
	url, err := buildURL(tokens)
	if err != nil {
		term.printf("Couldn't build URL: %v\n", err)
		return
	}

	request, err := http.NewRequest(c.method, url.String(), nil)
	if err != nil {
		term.printf("Couldn't build request: %v\n", err)
		return
	}

	//
	// User-specified params.
	//
	params := request.URL.Query()
	for k, v := range settings.Params {
		params.Add(k, v)
	}
	request.URL.RawQuery = params.Encode()

	for k, v := range settings.Headers {
		request.Header[k] = []string{v}
	}

	err = doRequest(term, request)
	if err != nil {
		term.printf("Error performing %s: %v\n", c.method, err)
	}
}

func buildURL(tokens []string) (*url.URL, error) {
	root := settings.Settings["root"]
	rootURL, _ := url.Parse("")
	tokenURL, _ := url.Parse("")

	var err error

	if len(root) > 0 {
		rootURL, err = url.Parse(root)
		if err != nil {
			return nil, err
		}
	}

	//
	// Any parameters that begin with a quote or @ cannot be the URL.
	//
	for _, val := range tokens[1:] {
		if !strings.HasPrefix(val, "@") && !strings.HasPrefix(val, "#") {
			var err error
			tokenURL, err = url.Parse(val)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	// if len(tokens) > 1 {
	// 	var err error
	// 	tokenURL, err = url.Parse(tokens[1])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	url := rootURL.ResolveReference(tokenURL)
	if len(url.String()) == 0 {
		return nil, fmt.Errorf("No URL specified!")
	}

	return url, nil
}
