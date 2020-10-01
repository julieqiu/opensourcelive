package main

import "io/ioutil"

type Page struct {
	Body string
}

func main() {
	loadPage("static/turkeychili.md")
}

func loadPage(title string) (*Page, error) {
	filename := title + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Body: string(body)}, nil
}
