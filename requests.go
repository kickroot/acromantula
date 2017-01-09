package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

var transport = &http.Transport{DisableKeepAlives: false}
var client = &http.Client{Timeout: time.Second * 10, Transport: transport}

//
// doRequest takes the supplied Request object and attempts to
// execute it, displaying the response contents and possibly
// returning an error condition if one occured.
//
func doRequest(term *Term, req *http.Request) error {
	term.writeString("\n<<  ")
	term.underscore()
	term.printf("%v %v\n", req.Method, req.URL)
	term.reset()
	printHeaders(" > ", term, req.Header)

	response, err := client.Do(req)
	transport.CloseIdleConnections()
	if err != nil {
		return err
	}

	defer response.Body.Close()
	term.writeString("\n<<  ")
	term.underscore()
	term.printf("HTTP %v\n", response.Status)
	term.reset()
	printHeaders(" < ", term, response.Header)
	printResponse(term, response)
	return nil
}

func printHeaders(prompt string, term *Term, headers http.Header) {
	for header, values := range headers {
		if strings.ToLower(header) == "authorization" {
			term.printf("%v %v : [****************]\n", prompt, header)
		} else {
			term.printf("%v %v : %v\n", prompt, header, values)
		}
	}
}

func printResponse(term *Term, response *http.Response) {
	term.writeString("\n<<  ")
	term.underscore()
	term.writeString("Content:\n")
	term.reset()
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
