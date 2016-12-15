package main

var default_headers map[string]string =  make(map[string]string)
var user_headers map[string]string =  make(map[string]string)

func init() {
	build_default_headers()
}

func build_default_headers() {
  default_headers["Accept"] = "application/json"
  default_headers["Accept-Charset"] = "utf-8"
  default_headers["User-Agent"] ="Acromantula CLI 0.1.0"
}