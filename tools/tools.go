package tools

import (
	"os"
	"log"
	"io/ioutil"
)

func GetFileContent(path string) (content string) {
	pwd, err := os.Getwd()

	if err != nil {
		log.Println(err)
		return content
	}

	bytes, err := ioutil.ReadFile(pwd + "/" + path)

	if err != nil {
		log.Println(err)
		return content
	}

	content = string(bytes)

	return content
}