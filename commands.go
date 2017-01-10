package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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

	//
	// Enforcing the preconditions:  Either there must be no parameters after the GET/HEAD/delete
	// command, or the first parameter must be a relative URL.
	//
	if len(tokens) > 2 || strings.HasPrefix(tokens[1], "@") {
		term.printf("Usage: %s [URL path]\n", c.method)
		return
	}

	var urlToken string
	if len(tokens) == 2 {
		urlToken = tokens[1]
	}

	url, err := parseURL(config.settings.Settings["root"], urlToken)
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
	for k, v := range config.settings.Params {
		params.Add(k, v)
	}
	request.URL.RawQuery = params.Encode()

	for k, v := range config.settings.Headers {
		request.Header[k] = []string{v}
	}

	err = doRequest(term, request)
	if err != nil {
		term.printf("Error performing %s: %v\n", c.method, err)
	}
}

func parseURL(root, token string) (*url.URL, error) {

	if len(root) == 0 && len(token) == 0 {
		return nil, fmt.Errorf("Root and passed URL cannot both be empty")
	}

	rootURL, _ := url.Parse(root)
	tokenURL, _ := url.Parse(token)

	//
	// Absolute URLs don't utilize root at all
	//
	if tokenURL.IsAbs() || len(root) == 0 {
		return tokenURL, nil
	}

	if len(token) == 0 {
		return rootURL, nil
	}

	return rootURL.ResolveReference(tokenURL), nil
}

type httpBodyCommand struct {
	method string
}

func (c *httpBodyCommand) exec(tokens []string, term *Term, config *configuration) {

	//
	// Enforcing the preconditions:  Either there must be no parameters after the GET/HEAD/delete
	// command, or the first parameter must be a relative URL.
	//
	if len(tokens) > 3 || (len(tokens) == 2 && strings.HasPrefix(tokens[1], "@")) {
		term.printf("Usage: %s [<URL path> [@/path/to/data]]\n", c.method)
		return
	}

	var urlToken string
	if len(tokens) > 1 {
		urlToken = tokens[1]
	}

	postURL, err := parseURL(config.settings.Settings["root"], urlToken)
	if err != nil {
		term.printf("Couldn't build URL: %v\n", err)
		return
	}

	// Optional request body, may be either parameter or data based.
	var body []byte
	contentType := ""

	//
	// User-specified params, this is overridden by any explicitly set POST
	// data (see @ token)
	//
	params := url.Values{}
	for k, v := range config.settings.Params {
		params.Add(k, v)
	}
	if len(params) > 0 {
		contentType = "application/x-www-form-urlencoded"
		body = []byte(params.Encode())
	}

	//
	// If any of the tokens starts with @, this is the file path to a data file that should be posted.
	//
	for _, token := range tokens {
		if strings.HasPrefix(token, "@") {
			dataFile := strings.TrimPrefix(token, "@")
			data, err := ioutil.ReadFile(dataFile)
			if err != nil {
				term.printf("Could perform %s, cannot read %v: %v\n", c.method, dataFile, err)
				return
			}
			body = data
			contentType = contentTypes[strings.TrimPrefix(filepath.Ext(dataFile), ".")]
			break
		}
	}

	request, err := http.NewRequest(c.method, postURL.String(), bytes.NewReader(body))
	if err != nil {
		term.printf("Couldn't build request: %v\n", err)
		return
	}

	for k, v := range config.settings.Headers {
		request.Header[k] = []string{v}
	}

	// If no custom Content-Type has been specified, use what we've discovered
	if len(config.settings.Headers["Content-Type"]) == 0 {
		request.Header["Content-Type"] = []string{contentType}
	}

	err = doRequest(term, request)
	if err != nil {
		term.printf("Error performing %s: %v\n", c.method, err)
	}
}
