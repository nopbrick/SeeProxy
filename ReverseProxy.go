package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(teamServer string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(teamServer)
	if err != nil {
		log.Fatal("Error parsing the teamserver URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(request *http.Request) {
		originalDirector(request)
		ModifyRequest(request)
	}

	proxy.ErrorHandler = ErrorHandler()
	return proxy, nil
}

func ModifyRequest(request *http.Request) {
	request.Header.Set("X-Proxy", "Simple Proxy")
	ValidateRequest(request)
}

func ValidateRequest(request *http.Request) {
	profile := ParseProfile("lambda.profile")
	headers := request.Header

	for _, parameter := range profile {
		for key, _ := range parameter {
			_, ok := headers[key]
			if !ok {
				ctx, cancel := context.WithCancel(context.Background())
				request, _ = http.NewRequestWithContext(ctx, "OPTIONS", "http://localhost:1", nil)
				cancel()
				// TODO: cancel request when not compliant with profile
			}
		}
	}
}

func ErrorHandler() func(w http.ResponseWriter, r *http.Request, e error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	proxy, err := NewProxy("http://127.0.0.1:8000")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":8888", nil)
}
