package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	filename := "turkeychili.md"
	s, err := loadPage(filename)
	if err != nil {
		log.Fatalf("unable to load page: %s", err)
	}
	fmt.Println(s)
}

func loadPage(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
