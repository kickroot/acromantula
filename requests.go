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
	term.writeString(fmt.Sprintf("\n>>  Performing %v on %v\n", req.Method, req.URL))
	term.writeString(fmt.Sprintf("Content-length: %v\n", req.ContentLength))
	printHeaders(" > ", term, req.Header)

	response, err := client.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	term.writeString(fmt.Sprintf("\n<<  Server returned HTTP %v\n", response.Status))
	printHeaders(" < ", term, response.Header)
	printResponse(term, response)
	return nil
}

func printHeaders(prompt string, term *Term, headers http.Header) {
	// term.writeString(fmt.Sprintf("%v Headers\n", title))
	for header, values := range headers {
		if strings.ToLower(header) == "authorization" {
			term.writeString(fmt.Sprintf("%v %v : [****************]\n", prompt, header))
		} else {
			term.writeString(fmt.Sprintf("%v %v : %v\n", prompt, header, values))
		}
	}
}

func printResponse(term *Term, response *http.Response) {
	term.writeString(fmt.Sprintf("\n<<  Response content: %v bytes\n", response.ContentLength))
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
