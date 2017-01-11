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
