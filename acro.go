package main

import (
    "strings"
    "bufio"
    "fmt"
    "os"
    "net/url"
    "log"
)

var default_headers map[string]string =  make(map[string]string)
var user_headers map[string]string =  make(map[string]string)
var base url.URL

func main() {

  build_default_headers();
  set_base("http://localhost");
 

  for {
    input := read()    
    tokens := strings.Split(input, " ")
    if (len(strings.TrimSpace(tokens[0])) == 0) {
      continue;
    }

    switch strings.TrimSpace(tokens[0]) {
    case "exit" :
      os.Exit(0)
    case "headers" :
      handle_headers(tokens)
    case "base" :
      handle_base(tokens)
    default :
      fmt.Printf(input)
    }
  }    
}

func set_base(str string) {
  default_url, err := url.Parse(str)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }    

  base = *default_url;
}

func handle_base(tokens []string) {

  if (len(tokens) == 1) {
    fmt.Printf("%v\n", base.String())
  } else {
    set_base(strings.TrimSpace(tokens[1]))
  }
}

func handle_headers(tokens []string) {
  
  //
  // In the case of just 'headers' 
  //
  if (len(tokens) == 1) {
    fmt.Printf("Headers\n")
    for k,v := range all_headers() {
      fmt.Printf("%v : %v\n", k, v)
    }
    return
  }

  switch tokens[1] {
    case "set":
      if (len(tokens) < 4) {
        fmt.Printf("%v is missing a value\n");
      } else {
        user_headers[strings.TrimSpace(tokens[2])] = strings.TrimSpace(tokens[3])
      }
    default :
      fmt.Printf("Unknown option %v\n", tokens[1])
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

func build_default_headers() {
  default_headers["Accept"] = "application/json"
  default_headers["Accept-Charset"] = "utf-8"
  default_headers["User-Agent"] ="Acromantula CLI 0.1.0"
}

func read() string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("acro >> ")
  text, _ := reader.ReadString('\n')
  return text
}
