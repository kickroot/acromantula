package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var transport = &http.Transport{DisableKeepAlives: false}
var client = &http.Client{Timeout: time.Second * 10, Transport: transport}

type httpCommand struct {
	method string
}

func (c *httpCommand) description() string {
	return fmt.Sprintf("Executes an HTTP %s request", c.method)
}

func (c *httpCommand) usage() string {
	return fmt.Sprintf("[url]")
}

func (c *httpCommand) exec(tokens []string, term *Term, config *configuration) {

	//
	// Enforcing the preconditions:  Either there must be no parameters after the GET/HEAD/delete
	// command, or the first parameter must be a relative URL.
	//
	if len(tokens) > 2 || strings.HasPrefix(tokens[1], "@") {
		term.printf("Usage: %s [URL path]\n", c.method)
		return
	}

	var urlToken string
	if len(tokens) == 2 {
		urlToken = tokens[1]
	}

	url, err := parseURL(config.settings.Settings["root"], urlToken)
	if err != nil {
		term.printf("Couldn't build URL: %v\n", err)
		return
	}

	request, err := http.NewRequest(c.method, url.String(), nil)
	if err != nil {
		term.printf("Couldn't build request: %v\n", err)
		return
	}

	//
	// User-specified params.
	//
	params := request.URL.Query()
	for k, v := range config.settings.Params {
		params.Add(k, v)
	}
	request.URL.RawQuery = params.Encode()

	for k, v := range config.settings.Headers {
		request.Header[k] = []string{v}
	}

	err = doRequest(term, request)
	if err != nil {
		term.printf("Error performing %s: %v\n", c.method, err)
	}
}

func parseURL(root, token string) (*url.URL, error) {

	if len(root) == 0 && len(token) == 0 {
		return nil, fmt.Errorf("Root and passed URL cannot both be empty")
	}

	rootURL, _ := url.Parse(root)
	tokenURL, _ := url.Parse(token)

	//
	// Absolute URLs don't utilize root at all
	//
	if tokenURL.IsAbs() || len(root) == 0 {
		return tokenURL, nil
	}

	if len(token) == 0 {
		return rootURL, nil
	}

	return rootURL.ResolveReference(tokenURL), nil
}

type httpBodyCommand struct {
	method string
}

func (c *httpBodyCommand) usage() string {
	return fmt.Sprintf("[<url> [@/path/to/file]]")
}

func (c *httpBodyCommand) description() string {
	return fmt.Sprintf("Executes an HTTP %s request", c.method)
}

func (c *httpBodyCommand) exec(tokens []string, term *Term, config *configuration) {

	//
	// Enforcing the preconditions:  Either there must be no parameters after the GET/HEAD/delete
	// command, or the first parameter must be a relative URL.
	//
	if len(tokens) > 3 || (len(tokens) == 2 && strings.HasPrefix(tokens[1], "@")) {
		term.printf("Usage: %s [<URL path> [@/path/to/data]]\n", c.method)
		return
	}

	var urlToken string
	if len(tokens) > 1 {
		urlToken = tokens[1]
	}

	postURL, err := parseURL(config.settings.Settings["root"], urlToken)
	if err != nil {
		term.printf("Couldn't build URL: %v\n", err)
		return
	}

	// Optional request body, may be either parameter or data based.
	var body []byte
	contentType := ""

	//
	// User-specified params, this is overridden by any explicitly set POST
	// data (see @ token)
	//
	params := url.Values{}
	for k, v := range config.settings.Params {
		params.Add(k, v)
	}
	if len(params) > 0 {
		contentType = "application/x-www-form-urlencoded"
		body = []byte(params.Encode())
	}

	//
	// If any of the tokens starts with @, this is the file path to a data file that should be posted.
	//
	for _, token := range tokens {
		if strings.HasPrefix(token, "@") {
			dataFile := strings.TrimPrefix(token, "@")
			data, err := ioutil.ReadFile(dataFile)
			if err != nil {
				term.printf("Could perform %s, cannot read %v: %v\n", c.method, dataFile, err)
				return
			}
			body = data
			contentType = contentTypes[strings.TrimPrefix(filepath.Ext(dataFile), ".")]
			break
		}
	}

	request, err := http.NewRequest(c.method, postURL.String(), bytes.NewReader(body))
	if err != nil {
		term.printf("Couldn't build request: %v\n", err)
		return
	}

	for k, v := range config.settings.Headers {
		request.Header[k] = []string{v}
	}

	// If no custom Content-Type has been specified, use what we've discovered
	if len(config.settings.Headers["Content-Type"]) == 0 {
		request.Header["Content-Type"] = []string{contentType}
	}

	err = doRequest(term, request)
	if err != nil {
		term.printf("Error performing %s: %v\n", c.method, err)
	}
}

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
