package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	settings.Headers["User-Agent"] = "Acromantula CLI 0.1.0"

	settings.Settings["root"] = "http://localhost"

	return &configuration{name: "acro", path: "", settings: *settings}
}
