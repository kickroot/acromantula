package main

import (
    "net/http"
    "encoding/json"
    "bytes"
    "log"
    "fmt"
)

var client = &http.Client{}


func performGet(term *Term, url string, headers map[string]string) {

  response, err := client.Get(url)
  if (err != nil) {
    fmt.Printf("Couldn't perform GET: %v\r\n", err)
  } else {
    defer response.Body.Close()
    fmt.Printf("Server returned %v\r\n", response.Status)

    for header,values := range response.Header {
      for value := range values {
        fmt.Printf("%v: %v\r\n", header, value)  
      }    
    }

    if (response.ContentLength != 0) {
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