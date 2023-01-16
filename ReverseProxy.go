package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(teamServer string, profile []ProfileParameters) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(teamServer)
	if err != nil {
		log.Fatal("Error parsing the teamserver URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(request *http.Request) {
		if !ValidateRequest(request, profile) {
			ctx, cancel := context.WithCancel(request.Context())
			*request = *request.WithContext(ctx)
			cancel()
			println("Terminated request from:", request.RemoteAddr, "(reason: request not compliant) ")
		} else {
			originalDirector(request)
		}

	}

	proxy.ErrorHandler = ErrorHandler()
	return proxy, nil
}

func ValidateRequest(request *http.Request, profile []ProfileParameters) bool {
	headers := request.Header

	for _, parameter := range profile {
		for key, _ := range parameter {
			_, ok := headers[key]
			if !ok {
				return false
			}
		}
	}
	return true
}

func ErrorHandler() func(w http.ResponseWriter, r *http.Request, e error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		//fmt.Printf("Terminated request from: %v \n", req.RemoteAddr)
		return
	}
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	teamserver := flag.String("teamserver", "127.0.0.1:8080", "Teamserver in format <IP>:<PORT>")
	profileFile := flag.String("profile", "my.profile", "Path to malleable profile")
	flag.Parse()
	profile := ParseProfile(*profileFile)

	proxy, err := NewProxy("http://"+*teamserver, profile)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":8888", nil)
}
