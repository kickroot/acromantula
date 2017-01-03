package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os/user"
	"path/filepath"
	"strings"
)

var root string
var term *Term
var settings *Settings

func main() {
	term = createTerm(0)

	usr, err := user.Current()
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't get the current user: %s", err))
	} else if len(usr.HomeDir) == 0 {
		term.writeString(fmt.Sprintf("Couldn't find home folder for %v", usr))
	}

	settingsFile := filepath.Join(usr.HomeDir, ".acromantula", "settings.yml")
	settings, err = initSettings(settingsFile)
	if err != nil {
		term.writeString(fmt.Sprintf("Error loading settings from %s, : %s", settingsFile, err))
	}
	term.setPrompt(settings.Settings["prompt"] + " >> ")
	term.writeString("Hit Ctrl+D to quit\n")

	defer term.restoreTerm()

	for {
		tokens, err := term.readline()

		if err != nil {
			term.writeString(fmt.Sprintf("\nExiting....%v\n", err))
			break
		}

		if len(tokens) == 0 || len(tokens[0]) == 0 {
			continue
		}

		switch tokens[0] {
		case "header":
			handleHeaders(tokens)
		case "headers":
			handleHeaders(tokens)
		case "set":
			handleSet(tokens)
		case "settings":
			handleSet(tokens)
		case "params":
			handleParams(tokens)
		case "get":
			handleHTTPGet(tokens)
		case "post":
			handleHTTPPost(tokens)
		default:
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		}
	}
}

func setRoot(str string) {
	root = str
}

func handleParams(tokens []string) {

	if len(tokens) == 1 {
		for key, value := range settings.Params {
			term.writeString(fmt.Sprintf("  %v => %v\n", key, value))
		}
		return
	}

	switch strings.ToLower(tokens[1]) {
	case "clear":
		if len(tokens) != 3 {
			term.writeString("Try 'params clear <key>'\n")
		} else {
			delete(settings.Params, tokens[2])
		}
	case "set":
		if len(tokens) != 4 {
			term.writeString("Try 'params set <key> <value>'\n")
		} else {
			settings.Params[tokens[2]] = tokens[3]
		}
	default:
		term.writeString(fmt.Sprintf("Unknown option '%v', try 'set' or 'clear'\n", tokens[1]))
	}
}

func handleSet(tokens []string) {

	if len(tokens) == 1 {
		for key, value := range settings.Settings {
			term.writeString(fmt.Sprintf("  %v => %v\n", key, value))
		}
		return
	}

	if len(tokens) < 3 {
		term.writeString(fmt.Sprintf("%v needs a value\n", tokens[1]))
	} else {
		settings.Settings[tokens[1]] = tokens[2]
		term.setPrompt(settings.Settings["prompt"] + " >> ")
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

	if len(tokens) > 1 {
		var err error
		tokenURL, err = url.Parse(tokens[1])
		if err != nil {
			return nil, err
		}
	}

	url := rootURL.ResolveReference(tokenURL)
	if len(url.String()) == 0 {
		return nil, fmt.Errorf("No URL specified!")
	}

	return url, nil
}

func handleHTTPPost(tokens []string) {
	postUrl, err := buildURL(tokens)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build URL: %v\n", err))
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
	for k, v := range settings.Params {
		params.Add(k, v)
	}
	if len(params) > 0 {
		contentType = "application/x-www-form-urlencoded"
		body = []byte(params.Encode())
	}

	request, err := http.NewRequest("POST", postUrl.String(), bytes.NewReader(body))
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build request: %v\n", err))
		return
	}

	for k, v := range settings.Headers {
		request.Header[k] = []string{v}
	}
	if len(contentType) != 0 {
		request.Header["Content-Type"] = []string{contentType}
	}

	err = doRequest(term, request)
	if err != nil {
		term.writeString(fmt.Sprintf("Error performing POST: %v\n", err))
	}
}

func handleHTTPGet(tokens []string) {
	url, err := buildURL(tokens)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build URL: %v\n", err))
		return
	}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build request: %v\n", err))
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
		term.writeString(fmt.Sprintf("Error performing GET: %v\n", err))
	}
}

func handleHeaders(tokens []string) {

	//
	// In the case of just 'headers'
	//
	if len(tokens) == 1 {
		term.writeString("Headers\n")
		for k, v := range settings.Headers {
			term.writeString(fmt.Sprintf("%v => %v\n", k, v))

		}
		return
	}

	switch tokens[1] {
	case "set":
		if len(tokens) < 4 {
			term.writeString(fmt.Sprintf("%v is missing a value\r\n", tokens[2]))
		} else {
			settings.Headers[tokens[2]] = tokens[3]
		}
	default:
		term.writeString(fmt.Sprintf("Unknown option %v\r\n", tokens[1]))
	}

}
