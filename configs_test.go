package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()
	if config.name != "acro" {
		t.Fatalf("Incorrect default name: %v", config.name)
	}
}

func TestWriteConfigToEmptyPath(t *testing.T) {
	config := defaultConfig()
	err := config.writeConfig()
	if err == nil {
		t.Fatalf("Err should be non-nil!")
	}
}

func TestWriteThenReadConfig(t *testing.T) {
	file, _ := ioutil.TempFile("", "")
	defer os.Remove(file.Name())

	config := defaultConfig()
	config.name = "test"
	config.path = file.Name()
	config.settings.Headers["_test"] = "_test_value"

	err := config.writeConfig()
	if err != nil {
		t.Fatalf("Error on writing configuration: %v", err)
	}

	config2, err := loadConfig("test2", file.Name())
	if err != nil {
		t.Fatalf("Error on reading configuration: %v", err)
	}

	if config2.name != "test2" {
		t.Fatalf("Incorrect config name, expected [test2] but found [%v]", config2.name)
	}

	for k, v := range config.settings.Headers {
		if config2.settings.Headers[k] != v {
			t.Fatalf("Bad header value for %v, expected [%v] but found [%v]", k, v, config2.settings.Headers[k])
		}
	}
}
