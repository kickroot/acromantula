package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var client = &http.Client{}

//
// doRequest takes the supplied Request object and attempts to
// execute it, displaying the response contents and possibly
// returning an error condition if one occured.
//
func doRequest(term *Term, req *http.Request) error {

	printHeaders("Request", term, req.Header)

	response, err := client.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	term.writeString(fmt.Sprintf("\nServer returned %v\n", response.Status))
	printHeaders("Response", term, response.Header)
	printResponse(term, response)
	return nil
}

func printHeaders(title string, term *Term, headers http.Header) {
	term.writeString(fmt.Sprintf("%v Headers\n", title))
	for header, values := range headers {
		term.writeString(fmt.Sprintf(" %v <= %v\n", header, values))
	}
}

func printResponse(term *Term, response *http.Response) {
	term.writeString(fmt.Sprintf("Response content: %v bytes\n", response.ContentLength))
	if response.ContentLength != 0 {
		formatted := new(bytes.Buffer)
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		bytes := buf.Bytes()
		error := json.Indent(formatted, bytes, "", "  ")
		if error != nil {
			term.writeBytes(bytes)
		} else {
			term.writeBytes(formatted.Bytes())
		}
		term.writeString("\n")
	}
}

func performGet(term *Term, url string, settings *Settings) {

	request, err := http.NewRequest("GET", url, nil)

	for k, v := range settings.Headers {
		request.Header[k] = []string{v}
	}

	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't create GET request: %v\n", err))
		return
	}

	term.writeString(fmt.Sprintf("\nPerforming GET on %v\n", url))
	term.writeString("Request Headers\n")

	for k, v := range request.Header {
		if k == "Authorization" {
			v = []string{"****************"}
		}
		term.writeString(fmt.Sprintf(" %v => %v\n", k, v))

	}

	response, err := client.Do(request)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't perform GET request: %v\n", err))
	} else {
		defer response.Body.Close()
		term.writeString(fmt.Sprintf("\nServer returned %v\n", response.Status))

		term.writeString("Response Headers\n")
		for header, values := range response.Header {
			term.writeString(fmt.Sprintf(" %v <= %v\n", header, values))
		}

		if response.ContentLength != 0 {
			formatted := new(bytes.Buffer)
			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			bytes := buf.Bytes()
			error := json.Indent(formatted, bytes, "", "  ")
			if error != nil {
				term.writeBytes(bytes)
			} else {
				term.writeBytes(formatted.Bytes())
			}
			term.writeString("\n")
		}
	}
}
