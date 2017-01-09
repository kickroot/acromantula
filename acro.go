package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const defaultConfigName string = "acro"

var configRoot string
var root string
var term *Term
var settings *Settings
var config *configuration

var headersCommand *mapCommand
var paramsCommand *mapCommand
var settingsCommand *mapCommand
var getCommand *httpCommand
var deleteCommand *httpCommand
var headCommand *httpCommand
var configCommand *configurationCommand

// The name of the currently applied configuration
var currentConfig = defaultConfigName

func main() {
	var err error

	term = createTerm(0)
	config = defaultConfig()

	//
	// All configurations are loaded/saved relative to the config root.  By default this is ~/.acromantula but
	// can be changed via the ACRO_CONFIG_ROOT env var.
	//
	configRoot, err = getConfigRoot()
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't determine config root, using defaults: %s\n", err))
	} else {

		//
		// Determine the path of the config file we are starting with.  This call cannot fail if the configRoot
		// isn't empty, so we can ignore the error value.
		//
		configFile, _ := getConfigPath(currentConfig)

		//
		// Attempt to read the default configuration, if it doesn't exist create it.
		//
		conf, err := loadConfig(currentConfig, configFile)
		if os.IsNotExist(err) {
			term.writeString("No settings file found, using defaults\n")
			configFile, err = getConfigPath(currentConfig)
			if err != nil {
				term.writeString(fmt.Sprintf("%s\n", err))
			} else {
				config.path = configFile
				term.writeString(fmt.Sprintf("Writing default config to %s\n", configFile))
				err = config.writeConfig()
				if err != nil {
					term.writeString(fmt.Sprintf("Couldn't save default config: %s\n", err))
				}
			}
		} else if err != nil {
			term.writeString(fmt.Sprintf("Error loading config from %s, : %s\n", configFile, err))
		} else {
			config = conf
		}
	}

	if err != nil {
		term.writeString(fmt.Sprintf("Error determining settings file for default configurationm using defaults: %s\n", err))
	} else {

	}
	settings = &config.settings
	updatePrompt()
	term.writeString("Hit Ctrl+D to quit\n")

	defer term.restoreTerm()

	initCommands(config)

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
			headersCommand.exec(tokens, term, config)
		case "headers":
			headersCommand.exec(tokens, term, config)
		case "set":
			settingsCommand.exec(tokens, term, config)
			updatePrompt()
		case "settings":
			settingsCommand.exec(tokens, term, config)
			updatePrompt()
		case "param":
			paramsCommand.exec(tokens, term, config)
		case "params":
			paramsCommand.exec(tokens, term, config)
		case "get":
			getCommand.exec(tokens, term, config)
		case "delete":
			deleteCommand.exec(tokens, term, config)
		case "head":
			headCommand.exec(tokens, term, config)
		case "post":
			handleHTTPPost(tokens)
		case "configs":
			configCommand.exec(tokens, term, config)
		case "config":
			configCommand.exec(tokens, term, config)
		default:
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		}
	}
}

func initCommands(config *configuration) {
	headersCommand = &mapCommand{backingMap: config.settings.Headers}
	paramsCommand = &mapCommand{backingMap: config.settings.Params}
	settingsCommand = &mapCommand{backingMap: config.settings.Settings}
	getCommand = &httpCommand{method: "GET"}
	deleteCommand = &httpCommand{method: "DELETE"}
	headCommand = &httpCommand{method: "HEAD"}
	configCommand = &configurationCommand{}
}

func handleHTTPPost(tokens []string) {
	postURL, err := buildURL(tokens)
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

	//
	// If any of the tokens starts with @, this is the file path to a data file that should be posted.
	//
	for _, token := range tokens {
		if strings.HasPrefix(token, "@") {
			dataFile := strings.TrimPrefix(token, "@")
			data, err := ioutil.ReadFile(dataFile)
			if err != nil {
				term.writeString(fmt.Sprintf("Could perform POST, cannot read %v: %v\n", dataFile, err))
				return
			}
			body = data
			contentType = contentTypes[strings.TrimPrefix(filepath.Ext(dataFile), ".")]
			break
		}
	}

	request, err := http.NewRequest("POST", postURL.String(), bytes.NewReader(body))
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't build request: %v\n", err))
		return
	}

	for k, v := range settings.Headers {
		request.Header[k] = []string{v}
	}

	// If no custom Content-Type has been specified, use what we've discovered
	if len(settings.Headers["Content-Type"]) == 0 {
		request.Header["Content-Type"] = []string{contentType}
	}

	err = doRequest(term, request)
	if err != nil {
		term.writeString(fmt.Sprintf("Error performing POST: %v\n", err))
	}
}

func updatePrompt() {
	//
	// If the current config has a custom setting for 'prompt', us it, otherwise
	// use the config name
	//
	prompt := config.name

	if len(settings.Settings["prompt"]) > 0 {
		prompt = settings.Settings["prompt"]
	}

	term.setPrompt(prompt)
}
