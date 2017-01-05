package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const defaultConfigName string = "acro"

var configRoot string
var root string
var term *Term
var settings *Settings
var config *configuration

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
		case "configs":
			handleConfig(tokens)
		case "config":
			handleConfig(tokens)
		default:
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		}
	}
}

func handleConfig(tokens []string) {

	//
	// A 'config' by itself just prompts for the current configuration
	//
	if len(tokens) == 1 {
		term.writeString(fmt.Sprintf("Current config: %s [%s]\n", config.name, config.path))
		return
	}

	switch tokens[1] {
	case "save":
		// With only 2 params, save to the current config
		targetConfig := currentConfig
		if len(tokens) > 2 {
			targetConfig = tokens[2]
		}

		configFile, err := getConfigPath(targetConfig)
		if err != nil {
			term.writeString(fmt.Sprintf("Couldn't save %v: %v", targetConfig, err))
			return
		}

		config.name = targetConfig
		config.path = configFile
		term.writeString(fmt.Sprintf("Saving %v to %v\n", config.name, config.path))
		err = config.writeConfig()
		if err != nil {
			term.writeString(fmt.Sprintf("Couldn't save %v: %v", targetConfig, err))
			return
		}
		updatePrompt()
	case "list":
		printConfigs(configRoot)
	case "load":
		if len(tokens) < 3 {
			term.writeString("Please supply a configuration name as well, such as 'config load acro'\n")
		} else {
			configFile, err := getConfigPath(tokens[2])
			if err != nil {
				term.writeString(fmt.Sprintf("Couldn't load %v: %v", tokens[2], err))
				return
			}
			conf, err := loadConfig(tokens[2], configFile)
			if err != nil {
				term.writeString(fmt.Sprintf("Couldn't load %v: %v", tokens[2], err))
				return
			}

			config = conf
			settings = &conf.settings
			currentConfig = config.name
			updatePrompt()
		}
	default:
		term.writeString(fmt.Sprintf("Unknown option '%s', try one of [save, list, load]\n", tokens[1]))
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
		updatePrompt()
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
		for k, v := range settings.Headers {
			term.writeString(fmt.Sprintf(" %v => %v\n", k, v))

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

func getConfigRoot() (string, error) {

	root = os.Getenv("ACRO_CONFIG_ROOT")
	if len(root) > 0 {
		return root, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".acromantula"), nil
}

func getConfigPath(configName string) (string, error) {

	if len(configRoot) == 0 {
		return "", fmt.Errorf("Cannot determine config location because config root is not known.")
	}
	return filepath.Join(configRoot, configName+".yml"), nil
}

func printConfigs(configRoot string) {

	if len(configRoot) == 0 {
		term.writeString("Can't print configs, no config root defined\n")
	}

	files, err := ioutil.ReadDir(configRoot)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't list configurations: %s\n", err))
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			configName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			term.writeString(fmt.Sprintf(" %v\n", configName))
		}
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
