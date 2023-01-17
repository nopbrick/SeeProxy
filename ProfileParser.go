package main

import (
	"os"
	"regexp"
	"strings"
)

type ProfileHeaders map[string][]string

type ProfileUserAgent map[string]string

func ParseProfile(file string) ([]ProfileHeaders, ProfileUserAgent) {

	rHTTPGET := regexp.MustCompile(`(?m)^http.*(.*)\{(.|\n)*?^}`)
	rClientHTTP := regexp.MustCompile(`(?m)^.client.*(.*)\{(.|\n)*?}\n.}`)
	rHeader := regexp.MustCompile(`(?m)header*.*"`)
	rUserAgent := regexp.MustCompile(`(?m)set useragent*.*"`)
	fileContent, _ := os.ReadFile(file)
	var empty []string
	empty = append(empty, "")

	var headers []ProfileHeaders

	for _, match := range rHTTPGET.FindAll(fileContent, -1) { // find every http-VERB block

		clientBlock := rClientHTTP.Find(match)
		for _, head := range rHeader.FindAll(clientBlock, -1) { // find every header option in clientblock
			tmp := strings.Replace(string(head), `"`, "", -1) // trim quotes
			h := strings.SplitN(tmp, " ", 3)
			if len(h) > 2 {
				headers = append(headers, ProfileHeaders{h[1]: strings.Split(h[len(h)-1], " ")})
			} else {
				headers = append(headers, ProfileHeaders{h[1]: empty})
			}
		}
	}
	userAgentParameter := strings.Replace(string(rUserAgent.Find(fileContent)), `"`, "", -1)
	userAgent := strings.SplitN(userAgentParameter, " ", 3)
	if len(userAgent) > 1 {
		profileUserAgent := ProfileUserAgent{userAgent[1]: userAgent[2]}
		return headers, profileUserAgent
	}
	return headers, nil
}
