package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

const lenPath = len("/view/")

var templates = make(map[string]*template.Template)

var titleValidate = regexp.MustCompile("^[a-zA-z0-9]+$")

const expend_string = ".txt"

func init() {
	for _, templ := range []string{"edit", "view"} {
		t := template.Must(template.ParseFiles(templ + ".html"))
		templates[templ] = t
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
	title = r.URL.Path[lenPath:]
	if !titleValidate.MatchString(title) {
		http.NotFound(w, r)
		err = errors.New("Invalid Page")
		log.Print(err)
	}
	return
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		err = errors.New("ショテイノディレクトリナイヨ")
		log.Print(err)
		return
	}

	var path []string
	var fileName []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), expend_string) {
			fileName = strings.Split(string(file.Name()), expend_string)
			path = append(path, fileName[0])
		}
	}

	if path == nil {
		err = errors.New("テキストファイルナイヨ")
		log.Print(err)
	}
	t := template.Must(template.ParseFiles("top.html"))
	err = t.Execute(w, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	lenderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	lenderTemplate(w, "edit", p)
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
		title := r.URL.Path[lenPath:]
		if !titleValidate.MatchString(title) {
			http.NotFound(w, r)
			err := errors.New("Invalid Page")
			log.Print(err)
			return
		}
		fn(w, r, title)
	}
}

func lenderTemplate(w http.ResponseWriter, templ string, p *Page) {

	t, _ := templates[templ].ParseFiles(templ + ".html")
	err := t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, err
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/top/", topHandler)
	http.ListenAndServe(":8080", nil)
}
