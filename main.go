package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"html/template"

	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Body template.HTML
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	addr := ":8080"
	log.Printf("Listening on addr %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello! Type the recipe you are looking for into the URL.")
		return
	}
	title := strings.TrimPrefix(r.URL.Path, "/")
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("view").Parse(`<div>{{.Body}}</div>`))

	if err := tmpl.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	b := renderBlackfriday(body)
	if err != nil {
		return nil, err
	}
	return &Page{Body: string(b)}, nil
}

func renderBlackfriday(body []byte) template.HTML {
	md := blackfriday.Run(body)
	return template.HTML(md)
}
