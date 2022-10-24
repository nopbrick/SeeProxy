package SeeProxy

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Handler(request http.Request) (http.Response, error) {

	var url *url.URL
	var body []byte
	var err error
	var outboundHeaders map[string]string

	teamserver := os.Getenv("TEAMSERVER")
	client := http.Client{}

	// Set to allow invalid HTTPS certs on the back-end server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Build our request URL as received to pass onto CS
	foo := "https://" + teamserver + "/" + request.RequestURI
	url, err = url.Parse(foo)
	if err != nil {
		log.Print(err)
	}

	// Extract any provided query parameters
	if request.URL.Query() != nil {
		q := url.Query()
		for key, value := range request.URL.Query() {
			q.Set(key, value[0])
		}
		url.RawQuery = q.Encode()
	}

	log.Print("url raw query: " + url.RawQuery)

	req, err := http.NewRequest(request.Method, url.String(), strings.NewReader(string(body)))
	if err != nil {
		log.Fatalf("Error pushing request to TeamServer: %v", err)
	}

	for key, values := range request.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Forward the request to our TeamServer
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error forwarding request to TeamServer: %v", err)
	}

	// Parse the TS response headers
	outboundHeaders = map[string]string{}

	for key, value := range resp.Header {
		outboundHeaders[key] = value[0]
	}

	// Store the TS response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error receiving request from TeamServer")
	}
	return http.Response{StatusCode: }
}
