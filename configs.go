package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type configuration struct {
	name     string
	path     string
	settings Settings
}

// writeConfig will write out the configuration to the specified path, overwriting any existing file.
func (c *configuration) writeConfig() error {
	if len(c.path) == 0 {
		return fmt.Errorf("Cannot write a config to an empty path.")
	}

	os.MkdirAll(filepath.Dir(c.path), 0700)
	return c.settings.writeSettings(c.path)
}

func loadConfig(name, path string) (*configuration, error) {
	settings, err := loadSettings(path)
	if err != nil {
		return nil, err
	}
	return &configuration{name: name, path: path, settings: *settings}, nil
}

// defaultConfig create a default configuration, suitable for one-time use or
// writing out an initial configuration on a new system.  This will have an empty file path.
func defaultConfig() *configuration {

	settings := defaultSettings()
	settings.Headers["Accept"] = "application/json, application/xml, application/xhtml+xml;q=0.9, text/html;q=0.9"
	settings.Headers["Accept-Charset"] = "utf-8"
	settings.Settings["root"] = "http://localhost"
	settings.Settings["prompt"] = "acro"

	return &configuration{name: "acro", path: "", settings: *settings}
}

func getConfigRoot() (string, error) {

	root := os.Getenv("ACRO_CONFIG_ROOT")
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
		term.printf("Couldn't list configurations: %s\n", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			configName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			term.printf(" %v\n", configName)
		}
	}
}
