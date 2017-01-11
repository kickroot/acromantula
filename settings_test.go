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
