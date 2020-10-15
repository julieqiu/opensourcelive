// The opensourcelive command starts a web server to render markdown recipe
// files as HTML web pages.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"html/template"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Page represents an HTML page.
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
	body, err := ioutil.ReadFile(filepath.Join("static", filename))
	if err != nil {
		return nil, err
	}

	b, err := renderMarkdown(body)
	if err != nil {
		return nil, err
	}
	return &Page{Body: b}, nil
}

func renderMarkdown(input []byte) (template.HTML, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Table, emoji.Emoji),
		goldmark.WithRendererOptions(html.WithUnsafe()))
	reader := text.NewReader(input)
	doc := md.Parser().Parse(reader)
	var b strings.Builder
	if err := md.Renderer().Render(&b, input, doc); err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}
