package main

type configurationCommand struct{}

func (c *configurationCommand) exec(tokens []string, term *Term, config *configuration) {
	//
	// A 'config' by itself just prompts for the current configuration
	//
	if len(tokens) == 1 {
		term.printf("Current config: %s [%s]\n", config.name, config.path)
		return
	}

	switch tokens[1] {
	case "save":
		// With only 2 params, save to the current config
		targetConfig := config.name
		if len(tokens) > 2 {
			targetConfig = tokens[2]
		}

		configFile, err := getConfigPath(targetConfig)
		if err != nil {
			term.printf("Couldn't save %v: %v\n", targetConfig, err)
			return
		}

		config.name = targetConfig
		config.path = configFile
		updatePrompt()
		initCommands(config)
		term.printf("Saving %v to %v\n", config.name, config.path)
		err = config.writeConfig()
		if err != nil {
			term.printf("Couldn't save %v: %v\n", targetConfig, err)
			return
		}
		updatePrompt()
		initCommands(config)
	case "list":
		printConfigs(configRoot)
	case "load":
		if len(tokens) < 3 {
			term.writeString("Please supply a configuration name as well, such as 'config load acro'\n")
		} else {
			configFile, err := getConfigPath(tokens[2])
			if err != nil {
				term.printf("Couldn't load %v: %v\n", tokens[2], err)
				return
			}
			conf, err := loadConfig(tokens[2], configFile)
			if err != nil {
				term.printf("Couldn't load %v: %v\n", tokens[2], err)
				return
			}

			// Todo: We need to encapsulate better, these are defined in acro.go
			config.name = conf.name
			config.path = conf.path
			config.settings = conf.settings
			updatePrompt()
			initCommands(config)
		}
	default:
		term.printf("Unknown option '%s', try one of [save, list, load]\n", tokens[1])
	}
}
