package main

import (
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
	proxy, err := NewProxy("http://127.0.0.1:1337")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":8888", nil)
}
