package main

import (
	"os"
	"regexp"
	"strings"
)

type Profile struct {
	ProfileHeadersGET
	ProfileHeadersPOST
	ProfileUserAgent
	ProfileURIsGET
	ProfileURIsPOST
}

type ProfileHeader map[string][]string
type ProfileUserAgent map[string]string
type ProfileURIsGET []string
type ProfileURIsPOST []string
type ProfileHeadersGET []ProfileHeader
type ProfileHeadersPOST []ProfileHeader

func ParseProfile(file string) Profile {

	rHTTPGET := regexp.MustCompile(`(?m)^http-get.*(.*)\{(.|\n)*?^}`)
	rHTTPPOST := regexp.MustCompile(`(?m)^http-post.*(.*)\{(.|\n)*?^}`)
	rClientHTTP := regexp.MustCompile(`(?m)^.client.*(.*)\{(.|\n)*?}\n.}`)
	rHeader := regexp.MustCompile(`(?m)header*.*"`)
	rUserAgent := regexp.MustCompile(`(?m)set useragent*.*"`)
	rURIs := regexp.MustCompile(`(?m)set uri*.*"`)

	fileContent, _ := os.ReadFile(file)
	var empty []string
	empty = append(empty, "")

	var headersGET []ProfileHeader
	var headersPOST []ProfileHeader
	var profileURIsGet ProfileURIsGET
	var profileURIsPOST ProfileURIsPOST

	for _, match := range rHTTPGET.FindAll(fileContent, -1) { // find every http-get block
		URIs := rURIs.Find(match)
		uriParameter := strings.Replace(string(URIs), `"`, "", -1) // trim quotes
		getURIs := strings.SplitN(uriParameter, " ", 3)
		getURIs = strings.Split(getURIs[2], " ")
		profileURIsGet = getURIs

		clientBlock := rClientHTTP.Find(match)
		for _, head := range rHeader.FindAll(clientBlock, -1) { // find every header option in clientblock
			tmp := strings.Replace(string(head), `"`, "", -1) // trim quotes
			h := strings.SplitN(tmp, " ", 3)
			if len(h) > 2 {
				headersGET = append(headersGET, ProfileHeader{h[1]: strings.Split(h[len(h)-1], " ")})
			} else {
				headersGET = append(headersGET, ProfileHeader{h[1]: empty})
			}
		}
	}

	for _, match := range rHTTPPOST.FindAll(fileContent, -1) { // find every http-post block
		URIs := rURIs.Find(match)
		uriParameter := strings.Replace(string(URIs), `"`, "", -1) // trim quotes
		postURIs := strings.SplitN(uriParameter, " ", 3)
		postURIs = strings.Split(postURIs[2], " ")
		profileURIsPOST = postURIs
		clientBlock := rClientHTTP.Find(match)
		for _, head := range rHeader.FindAll(clientBlock, -1) { // find every header option in clientblock
			tmp := strings.Replace(string(head), `"`, "", -1) // trim quotes
			h := strings.SplitN(tmp, " ", 3)
			if len(h) > 2 {
				headersPOST = append(headersPOST, ProfileHeader{h[1]: strings.Split(h[len(h)-1], " ")})
			} else {
				headersPOST = append(headersPOST, ProfileHeader{h[1]: empty})
			}
		}
	}
	userAgentParameter := strings.Replace(string(rUserAgent.Find(fileContent)), `"`, "", -1)
	userAgent := strings.SplitN(userAgentParameter, " ", 3)
	if len(userAgent) > 1 {
		profileUserAgent := ProfileUserAgent{userAgent[1]: userAgent[2]}
		profile := Profile{
			ProfileHeadersGET:  headersGET,
			ProfileHeadersPOST: headersPOST,
			ProfileUserAgent:   profileUserAgent,
			ProfileURIsGET:     profileURIsGet,
			ProfileURIsPOST:    profileURIsPOST,
		}
		return profile
	}
	return Profile{headersGET, headersPOST, nil, profileURIsGet, profileURIsPOST}
}
