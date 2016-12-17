package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var client = &http.Client{}

func performGet(term *Term, url string, headers *Headers) {

	request, err := http.NewRequest("GET", url, nil)

	for k, v := range headers.all() {
		request.Header[k] = []string{v}
	}

	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't create GET request: %v", err))
		return
	}

	term.writeString(fmt.Sprintf("\nPerforming GET on %v\n", url))
	term.writeString("Request Headers\n")
	for k, v := range request.Header {
		term.writeString(fmt.Sprintf(" %v => %v\n", k, v))
	}

	response, err := client.Do(request)
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't perform GET request: %v", err))
	} else {
		defer response.Body.Close()
		term.writeString(fmt.Sprintf("\nServer returned %v\n", response.Status))

		term.writeString("Response Headers\n")
		for header, values := range response.Header {
			term.writeString(fmt.Sprintf(" %v <= %v\n", header, values))
			// for value := range values {
			// 	fmt.Printf("%v: %v\r\n", header, values)
			// }
		}

		if response.ContentLength != 0 {
			formatted := new(bytes.Buffer)
			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			bytes := buf.Bytes()
			error := json.Indent(formatted, bytes, "", "  ")
			if error != nil {
				log.Println("JSON parse error: ", error)
			} else {
				term.writeBytes(formatted.Bytes())
				// io.Copy(os.Stdout, formatted)
				fmt.Printf("\r\n")
			}
		}
	}
}
