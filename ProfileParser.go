package main

import (
	"os"
	"regexp"
	"strings"
)

type ProfileParameters map[string][]string

func ParseProfile(file string) []ProfileParameters {

	rHTTPGET := regexp.MustCompile(`(?m)^http.*(.*)\{(.|\n)*?^}`)
	rClientHTTP := regexp.MustCompile(`(?m)^.client.*(.*)\{(.|\n)*?}\n.}`)
	rHeader := regexp.MustCompile(`(?m)header*.*"`)
	fileContent, _ := os.ReadFile(file)
	var empty []string
	empty = append(empty, "")

	var headers []ProfileParameters

	for _, match := range rHTTPGET.FindAll(fileContent, -1) { // find every http-VERB block
		clientBlock := rClientHTTP.Find(match)
		for _, head := range rHeader.FindAll(clientBlock, -1) { // find every header option in clientblock
			tmp := strings.Replace(string(head), `"`, "", -1) // trim quotes
			h := strings.SplitN(tmp, " ", 3)
			if len(h) > 2 {
				headers = append(headers, ProfileParameters{h[1]: strings.Split(h[len(h)-1], " ")})
			} else {
				headers = append(headers, ProfileParameters{h[1]: empty})
			}
		}
	}
	return headers
}
