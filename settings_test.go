package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMissingFile(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get pwd: %v", err)
	}

	settingsFile := filepath.Join(pwd, "tests/settings/missing.yml")
	_, err = initSettings(settingsFile)
	if err == nil {
		t.Fatalf("Expected a non-nil error value!")
	}
}

func TestBasicSettings(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get pwd: %v", err)
	}

	settingsFile := filepath.Join(pwd, "tests/settings/basic.yml")
	settings, err := initSettings(settingsFile)
	if err != nil {
		t.Fatalf("Found non-nil error on /home/basic.yml: %v", err)
	}

	if settings.Settings["setting1"] != "one" {
		t.Fatalf("Expected 'one' but found %v", settings.Settings["settings1"])
	}
}

func TestInvalidSettings(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get pwd: %v", err)
	}

	settingsFile := filepath.Join(pwd, "tests/settings/invalid.yml")
	_, err = initSettings(settingsFile)
	if err == nil {
		t.Fatalf("Expected a non-nil error value!")
	}

}
