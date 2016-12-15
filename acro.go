package main

import (
    "strings"
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "log"
)

var root string

func main() {
  fmt.Printf("Hit Ctrl+D to quit\r\n")
  initTerm();
  defer restoreTerm()

  setRoot("http://localhost");
 

  for {
    tokens, err := readTerm()
    
    if err != nil {
      fmt.Printf("\r\nExiting....%v\r\n", err)
      break
    }

    if len(tokens) == 0 || len(tokens[0]) == 0 {
      continue
    }

    switch tokens[0] {
    case "headers" :
      handle_headers(tokens)
    case "root" :
      handleRoot(tokens)
    case "get" :
      handleGet(tokens)
    default :
      fmt.Printf("Unknown command, %v\r\n", tokens[0])
    }
  }    
}

func setRoot(str string) {
  root = str
}

func handleRoot(tokens []string) {

  if (len(tokens) == 1) {
    fmt.Printf("%v\r\n", root)
  } else {
    setRoot(strings.TrimSpace(tokens[1]))
  }
}

func handleGet(tokens []string) {

  url := root
  if (len(tokens) > 1) {
    url = root + tokens[1]
  } 

  response, err := http.Get(url)
  if (err != nil) {
    fmt.Printf("Couldn't perform GET: %v\r\n", err)
  } else {
    defer response.Body.Close()
    fmt.Printf("Server returned %v\r\n", response.Status)
    if (response.ContentLength != 0) {
      formatted := new(bytes.Buffer)
      buf := new(bytes.Buffer)
      buf.ReadFrom(response.Body)
      bytes := buf.Bytes()
      error := json.Indent(formatted, bytes, "", "  ")
      if error != nil {
        log.Println("JSON parse error: ", error)
      } else {
        writeTerm(formatted.Bytes())
        // io.Copy(os.Stdout, formatted)
        fmt.Printf("\r\n")        
      }
    }
  }
}

func handle_headers(tokens []string) {
  
  //
  // In the case of just 'headers' 
  //
  if (len(tokens) == 1) {
    fmt.Printf("Headers\r\n")
    for k,v := range all_headers() {
      fmt.Printf("%v : %v\r\n", k, v)
    }
    return
  }

  switch tokens[1] {
    case "set":
      if (len(tokens) < 4) {
        fmt.Printf("%v is missing a value\r\n");
      } else {
        user_headers[strings.TrimSpace(tokens[2])] = strings.TrimSpace(tokens[3])
      }
    default :
      fmt.Printf("Unknown option %v\r\n", tokens[1])
  }

}

func all_headers() map[string] string {
  all_headers := make(map[string]string)

  for k, v := range default_headers {
    all_headers[k] = v
  }

  for k, v := range user_headers {
    all_headers[k] = v
  }

  return all_headers
}
