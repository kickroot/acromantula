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

import "testing"

func TestBasicParse(t *testing.T) {
	root := "https://example.com/"
	token := "me"
	expected := "https://example.com/me"

	url, err := parseURL(root, token)
	if err != nil {
		t.Fatalf("Expected a nil error value!")
	}

	if url.String() != expected {
		t.Fatalf("Expected %s, found %s", expected, url)
	}
}

func TestDoubleSlashes(t *testing.T) {
	root := "https://example.com/"
	token := "/me"
	expected := "https://example.com/me"

	url, err := parseURL(root, token)
	if err != nil {
		t.Fatalf("Expected a nil error value!")
	}

	if url.String() != expected {
		t.Fatalf("Expected %s, found %s", expected, url)
	}
}

func TestNoSlashes(t *testing.T) {
	root := "https://example.com"
	token := "me"
	expected := "https://example.com/me"

	url, err := parseURL(root, token)
	if err != nil {
		t.Fatalf("Expected a nil error value!")
	}

	if url.String() != expected {
		t.Fatalf("Expected %s, found %s", expected, url)
	}
}

func TestAbsoluteURL(t *testing.T) {
	root := "https://example.com"
	token := "https://api.example.com"
	expected := "https://api.example.com"

	url, err := parseURL(root, token)
	if err != nil {
		t.Fatalf("Expected a nil error value!")
	}

	if url.String() != expected {
		t.Fatalf("Expected %s, found %s", expected, url)
	}
}

func TestNoURLs(t *testing.T) {
	root := ""
	token := ""

	_, err := parseURL(root, token)
	if err == nil {
		t.Fatalf("Expected a non-nil error value!")
	}
}
