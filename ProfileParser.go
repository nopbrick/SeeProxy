package main

import (
	"os"
	"regexp"
	"strings"
)

type Profile struct {
	ProfileHeaders
	ProfileUserAgent
	ProfileURIsGET
}

type ProfileHeaders map[string][]string
type ProfileUserAgent map[string]string
type ProfileURIsGET []string

func ParseProfile(file string) ([]ProfileHeaders, ProfileUserAgent) {

	rHTTPGET := regexp.MustCompile(`(?m)^http-get.*(.*)\{(.|\n)*?^}`)
	rHTTPPOST := regexp.MustCompile(`(?m)^http-post.*(.*)\{(.|\n)*?^}`)
	rClientHTTP := regexp.MustCompile(`(?m)^.client.*(.*)\{(.|\n)*?}\n.}`)
	rHeader := regexp.MustCompile(`(?m)header*.*"`)
	rUserAgent := regexp.MustCompile(`(?m)set useragent*.*"`)
	rURIs := regexp.MustCompile(`(?m)set uri*.*"`)

	fileContent, _ := os.ReadFile(file)
	var empty []string
	empty = append(empty, "")

	var headers []ProfileHeaders

	for _, match := range rHTTPGET.FindAll(fileContent, -1) { // find every http-get block
		URIs := rURIs.Find(match)
		uriParameter := strings.Replace(string(URIs), `"`, "", -1) // trim quotes
		getURIs := strings.SplitN(uriParameter, " ", 3)
		getURIs = strings.Split(getURIs[2], " ")
		profileURIsGet := ProfileURIsGET(getURIs)
		println(profileURIsGet)

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

	for _, match := range rHTTPPOST.FindAll(fileContent, -1) { // find every http-post block

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
