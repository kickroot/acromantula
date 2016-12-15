package main

import (
    "strings"
    "fmt"
    "net/url"
)

var base url.URL

func main() {
  fmt.Printf("Hit Ctrl+D to quit\r\n")
  initTerm();
  defer restoreTerm()

  set_base("http://localhost");
 

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
    case "base" :
      handle_base(tokens)
    default :
      fmt.Printf("Unknown command, %v\r\n", tokens[0])
    }
  }    
}

func set_base(str string) {
  default_url, err := url.Parse(str)
  if err != nil {
    fmt.Printf("Invalid url %v:%v\r\n", str, err)
  } else {
    base = *default_url;
  }    
}

func handle_base(tokens []string) {

  if (len(tokens) == 1) {
    fmt.Printf("%\r\n", base.String())
  } else {
    set_base(strings.TrimSpace(tokens[1]))
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
