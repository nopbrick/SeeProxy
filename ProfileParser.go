package main

import (
	"os"
	"regexp"
	"strings"
)

type ProfileParametersPOST map[string]string
type ProfileParametersGET map[string]string

type Profile struct {
	ProfileParametersPOST
	ProfileParametersGET
}

func ParseProfile(file string) {
	rHTTPGET := regexp.MustCompile(`(?m)^http.*(.*)\{(.|\n)*?^}`)
	rClientHTTP := regexp.MustCompile(`(?m)^.client.*(.*)\{(.|\n)*?}\n.}`)
	rHeader := regexp.MustCompile(`(?m)header*.*"`)
	fileContent, _ := os.ReadFile(file)

	for _, match := range rHTTPGET.FindAll(fileContent, -1) {
		clientBlock := rClientHTTP.Find(match)
		tmp := rHeader.FindAll(clientBlock, -1)
		param := strings.Replace(string(tmp[0]), `"`, "", -1)
		//param = strings.Split(string(tmp[0]), " ")
		println(param)
		//profileParameters := ProfileParametersGET{"fsd": "1", "bar": "2"}
		//fmt.Println(string(tmp[0]))
		//fmt.Println(clientBlock)
	}

}

func main() {
	ParseProfile("lambda.profile")
}
