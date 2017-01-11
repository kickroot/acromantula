package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

	if len(headers["User-Agent"]) == 0 {
		headers["User-Agent"] = []string{fmt.Sprintf("Acromantula %s", acroVersion)}
	}

	for _, key := range sortHeaders(headers) {
		if strings.ToLower(key) == "authorization" {
			term.printf("%v %v : [****************]\n", prompt, key)
		} else {
			term.printf("%v %v : %v\n", prompt, key, headers[key])
		}
	}
}

func sortHeaders(h http.Header) []string {
	keys := make([]string, len(h))
	i := 0
	for k := range h {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
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
