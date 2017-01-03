package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
		if strings.ToLower(header) == "authorization" {
			term.writeString(fmt.Sprintf(" %v <= [****************]\n", header))
		} else {
			term.writeString(fmt.Sprintf(" %v <= %v\n", header, values))
		}
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
