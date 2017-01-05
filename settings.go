package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Settings map[string]string `yaml:"settings"`
	Headers  map[string]string `yaml:"headers"`
	Params   map[string]string `yaml:"params"`
}

func defaultSettings() *Settings {
	settings := Settings{}
	settings.Settings = make(map[string]string)
	settings.Headers = make(map[string]string)
	settings.Params = make(map[string]string)

	return &settings
}

func initSettings(settingsFile string) (*Settings, error) {
	settings := defaultSettings()

	settings.Headers["Accept"] = "application/json"
	settings.Headers["Accept-Charset"] = "utf-8"
	settings.Headers["User-Agent"] = "Acromantula CLI 0.1.0"

	settings.Settings["root"] = "http://localhost"

	bytes, err := ioutil.ReadFile(settingsFile)

	if err != nil {
		return settings, fmt.Errorf("Couldn't read %v: %v", settingsFile, err)
	}

	err = yaml.Unmarshal(bytes, &settings)
	if err != nil {
		return settings, fmt.Errorf("Couldn't unmarshal yml: %v", err)
	}

	return settings, nil
}

func loadSettings(path string) (*Settings, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	settings := defaultSettings()
	err = yaml.Unmarshal(bytes, settings)
	return settings, err
}

func (s *Settings) writeSettings(settingsFile string) error {
	bytes, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(settingsFile, bytes, 0600)
}
