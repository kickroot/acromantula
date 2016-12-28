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

	// usr, err := user.Current()
	// if err != nil {
	// 	return &settings, err
	// } else if len(usr.HomeDir) == 0 {
	// 	return &settings, errors.New(Sprintf("Coulidn't find home folder for %v", usr))
	// }

	// settingsFile := filepath.Join(usr.HomeDir, ".acromantula", "settings.yml")
	//Printf("Reading file %v", settingsFile)
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
