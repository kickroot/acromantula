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
