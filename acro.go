package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os/user"
	"path/filepath"
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
		case "get":
			handleGet(tokens)
		default:
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		}
	}
}

func setRoot(str string) {
	root = str
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

func handleGet(tokens []string) {
	root := settings.Settings["root"]
	rootURL, _ := url.Parse("")
	tokenURL, _ := url.Parse("")

	if len(root) > 0 {
		var err error
		rootURL, err = url.Parse(root)
		if err != nil {
			term.writeString(fmt.Sprintf("Bad root URL specified: %v\n", err))
			return
		}
	}

	if len(tokens) > 1 {
		var err error
		tokenURL, err = url.Parse(tokens[1])
		if err != nil {
			term.writeString(fmt.Sprintf("Bad GET URL specified: %v\n", err))
			return
		}
	}

	url := rootURL.ResolveReference(tokenURL)
	if len(url.String()) == 0 {
		term.writeString("No URL specified!\n")
	}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build request: %v\n", err))
		return
	}

	for k, v := range settings.Headers {
		request.Header[k] = []string{v}
	}

	err = doRequest(term, request)
	if err != nil {
		term.writeString(fmt.Sprintf("Error performing GET: %v\n", err))
	}

	// performGet(term, url.String(), settings)
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
