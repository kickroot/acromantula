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

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const acroVersion = "0.1.0-alpha"
const defaultConfigName = "default"

var configRoot string
var term *Term
var config *configuration

var headersCommand *mapCommand
var paramsCommand *mapCommand
var settingsCommand *mapCommand

var commands map[string]command

var getCommand = &httpCommand{method: "GET"}
var deleteCommand = &httpCommand{method: "DELETE"}
var headCommand = &httpCommand{method: "HEAD"}
var putCommand = &httpBodyCommand{method: "PUT"}
var postCommand = &httpBodyCommand{method: "POST"}
var configCommand = &configurationCommand{}

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
	term.printf("Acromantula %s\n", acroVersion)
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

		cmd := commands[strings.ToLower(tokens[0])]
		if cmd == nil {
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		} else {
			cmd.exec(tokens, term, config)
		}
		// switch tokens[0] {
		// case "header":
		// 	headersCommand.exec(tokens, term, config)
		// case "headers":
		// 	headersCommand.exec(tokens, term, config)
		// case "set":
		// 	settingsCommand.exec(tokens, term, config)
		// case "settings":
		// 	settingsCommand.exec(tokens, term, config)
		// case "param":
		// 	paramsCommand.exec(tokens, term, config)
		// case "params":
		// 	paramsCommand.exec(tokens, term, config)
		// case "get":
		// 	getCommand.exec(tokens, term, config)
		// case "delete":
		// 	deleteCommand.exec(tokens, term, config)
		// case "head":
		// 	headCommand.exec(tokens, term, config)
		// case "post":
		// 	postCommand.exec(tokens, term, config)
		// case "put":
		// 	putCommand.exec(tokens, term, config)
		// case "configs":
		// 	configCommand.exec(tokens, term, config)
		// case "config":
		// 	configCommand.exec(tokens, term, config)
		// default:

		// }
		updatePrompt()
	}
}

func initCommands(config *configuration) {

	commands = make(map[string]command)
	commands["get"] = &httpCommand{method: "GET"}
	commands["head"] = &httpCommand{method: "HEAD"}
	commands["delete"] = &httpCommand{method: "DELETE"}
	commands["post"] = &httpBodyCommand{method: "POST"}
	commands["put"] = &httpBodyCommand{method: "PUT"}
	commands["config"] = &configurationCommand{}
	commands["help"] = &helpCommand{}

	updateCommands(config)
}

func updateCommands(config *configuration) {
	commands["headers"] = &mapCommand{desc: "Headers for all HTTP(S) requests", backingMap: config.settings.Headers}
	commands["header"] = commands["headers"]
	commands["params"] = &mapCommand{desc: "Request parameters for all HTTP(S) requests", backingMap: config.settings.Params}
	commands["param"] = commands["params"]
	commands["settings"] = &mapCommand{desc: "Application level settings and preferences", backingMap: config.settings.Settings}
	commands["setting"] = commands["settings"]
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

func sortKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

func sortCommands(m map[string]command) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}
