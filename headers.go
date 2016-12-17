package main

// Headers we send to the server
type Headers struct {
	defaultHeaders map[string]string
	userHeaders    map[string]string
}

func createHeaders() *Headers {
	h := new(Headers)
	h.defaultHeaders = make(map[string]string)
	h.userHeaders = make(map[string]string)
	h.defaults()
	return h
}

func (h *Headers) defaults() {
	h.defaultHeaders["Accept"] = "application/json"
	h.defaultHeaders["Accept-Charset"] = "utf-8"
	h.defaultHeaders["User-Agent"] = "Acromantula CLI 0.1.0"
}

func (h *Headers) add(key string, value string) {
	h.userHeaders[key] = value
}

func (h *Headers) all() map[string]string {
	allHeaders := make(map[string]string)

	for k, v := range h.defaultHeaders {
		allHeaders[k] = v
	}

	for k, v := range h.userHeaders {
		allHeaders[k] = v
	}

	return allHeaders
}
