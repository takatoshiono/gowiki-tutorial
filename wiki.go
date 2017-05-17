package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title    string
	Body     []byte
	HTMLBody template.HTML
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var linkPattern = regexp.MustCompile(`\[((?:[A-Z]{1}[a-z]+){1,})\]`)

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("data/"+filename, p.Body, 0600)
}

func (p *Page) convertBodyToHTML() {
	escapedHTML := template.HTMLEscapeString(string(p.Body))
	p.HTMLBody = template.HTML(linkPattern.ReplaceAllStringFunc(escapedHTML, convertLink))
}

func convertLink(s string) string {
	name := s[1 : len(s)-1]
	return fmt.Sprintf("<a href=\"/view/%s\">%s</a>", name, name)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	p.convertBodyToHTML()
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	port := flag.Int("p", 8080, "tcp port number")
	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("public/styles"))))

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
