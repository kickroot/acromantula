package main

import (
	"errors"
	. "fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Settings map[string]string `yaml:"settings"`
	Headers  map[string]string `yaml:"headers"`
}

func initSettings(settingsFile string) (*Settings, error) {
	settings := Settings{}
	settings.Settings = make(map[string]string)
	settings.Headers = make(map[string]string)

	settings.Headers["Accept"] = "application/json"
	settings.Headers["Accept-Charset"] = "utf-8"
	settings.Headers["User-Agent"] = "Acromantula CLI 0.1.0"

	settings.Settings["root"] = "http://localhost"
	settings.Settings["prompt"] = "acro"

	bytes, err := ioutil.ReadFile(settingsFile)

	if err != nil {
		return &settings, errors.New(Sprintf("Coulidn't read %v: %v", settingsFile, err))
	}

	err = yaml.Unmarshal(bytes, &settings)
	if err != nil {
		return &settings, errors.New(Sprintf("Coulidn't unmarshal yml: %v", err))
	}
	return &settings, nil
}
