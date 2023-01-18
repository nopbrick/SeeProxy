package main

import (
	"context"
	"flag"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(teamServer string, profile Profile) (*httputil.ReverseProxy, error) {
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

func ValidateRequest(request *http.Request, profile Profile) bool {
	headers := request.Header
	if profile.ProfileUserAgent != nil && request.UserAgent() != profile.ProfileUserAgent["useragent"] {
		return false
	}
	if request.Method == "GET" {
		if !slices.Contains(profile.ProfileURIsGET, request.RequestURI) {
			return false
		} else {
			for _, parameter := range profile.ProfileHeadersGET {
				for key, _ := range parameter {
					_, ok := headers[key]
					if !ok {
						return false
					}
				}
			}
		}
	}
	if request.Method == "POST" {
		if !slices.Contains(profile.ProfileURIsPOST, request.RequestURI) {
			return false
		} else {
			for _, parameter := range profile.ProfileHeadersPOST {
				for key, _ := range parameter {
					_, ok := headers[key]
					if !ok {
						return false
					}
				}
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
	teamserver := flag.String("teamserver", "127.0.0.1:8000", "Teamserver in format <IP>:<PORT>")
	profileFile := flag.String("profile", "lambda.profile", "Path to malleable profile")
	localPort := flag.String("port", "8080", "Local port to bind to")
	flag.Parse()

	profile := ParseProfile(*profileFile)

	proxy, err := NewProxy("http://"+*teamserver, profile)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":"+*localPort, nil)
}
