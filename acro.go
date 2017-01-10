package main

import (
	"fmt"
	"os"
)

const defaultConfigName string = "acro"

var configRoot string
var term *Term
var config *configuration

var headersCommand *mapCommand
var paramsCommand *mapCommand
var settingsCommand *mapCommand

var getCommand = &httpCommand{method: "GET"}
var deleteCommand = &httpCommand{method: "DELETE"}
var headCommand = &httpCommand{method: "HEAD"}
var putCommand = &httpBodyCommand{method: "PUT"}
var postCommand = &httpBodyCommand{method: "POST"}
var configCommand = &configurationCommand{}

// var configCommand *configurationCommand

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
		term.printf("Couldn't determine config root, using defaults: %s\n", err)
	} else {

		//
		// Determine the path of the config file we are starting with.  This call cannot fail if the configRoot
		// isn't empty, so we can ignore the error value.
		//
		configFile, _ := getConfigPath(defaultConfigName)

		//
		// Attempt to read the default configuration, if it doesn't exist create it.
		//
		conf, err := loadConfig(defaultConfigName, configFile)
		if os.IsNotExist(err) {
			term.writeString("No settings file found, using defaults\n")
			configFile, err = getConfigPath(defaultConfigName)
			if err != nil {
				term.printf("%s\n", err)
			} else {
				config.path = configFile
				term.printf("Writing default config to %s\n", configFile)
				err = config.writeConfig()
				if err != nil {
					term.printf("Couldn't save default config: %s\n", err)
				}
			}
		} else if err != nil {
			term.printf("Error loading config from %s, : %s\n", configFile, err)
		} else {
			config = conf
		}
	}

	if err != nil {
		term.printf("Error determining settings file for default configuration, using defaults: %s\n", err)
	} else {

	}

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
			postCommand.exec(tokens, term, config)
		case "put":
			putCommand.exec(tokens, term, config)
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
}

func updatePrompt() {
	//
	// If the current config has a custom setting for 'prompt', us it, otherwise
	// use the config name
	//
	prompt := config.name

	if len(config.settings.Settings["prompt"]) > 0 {
		prompt = config.settings.Settings["prompt"]
	}

	term.setPrompt(prompt)
}
