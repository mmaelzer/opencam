package templater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

var templates map[string]*template.Template
var subpages map[string]*template.Template
var partials map[string]*template.Template

func init() {
	helpers := template.FuncMap{
		"humanizeDate":   humanize.Time,
		"humanizeNumber": humanize.Comma,
		"json": func(v interface{}) template.JS {
			j, _ := json.Marshal(v)
			return template.JS(j)
		},
	}

	templates = make(map[string]*template.Template)
	partials = make(map[string]*template.Template)
	subpages = make(map[string]*template.Template)
	layout := "templater/layout.tmpl"

	pages, err := filepath.Glob("templater/pages/*.tmpl")
	if err != nil {
		panic(err)
	}

	includes, err := filepath.Glob("templater/includes/*.tmpl")
	if err != nil {
		panic(err)
	}

	subs, err := filepath.Glob("templater/subpages/*.tmpl")
	if err != nil {
		panic(err)
	}

	for _, partial := range includes {
		key := strings.Split(filepath.Base(partial), ".")[0]
		partials[key] = template.Must(template.New(partial).Funcs(helpers).ParseFiles(partial))
	}

	for _, subpage := range subs {
		files := append(includes, subpage)
		key := strings.Split(filepath.Base(subpage), ".")[0]
		subpages[key] = template.Must(template.New(subpage).Funcs(helpers).ParseFiles(files...))
	}

	for _, page := range pages {
		files := append(includes, page, string(layout))
		key := strings.Split(filepath.Base(page), ".")[0]
		templates[key] = template.Must(template.New(page).Funcs(helpers).ParseFiles(files...))
	}
}

func Build(name string, data interface{}) ([]byte, error) {
	tmpl, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("The template %s does not exist", name)
	}
	return execute(tmpl, "base", data)
}

func Partial(name string, data interface{}) ([]byte, error) {
	partial, ok := partials[name]
	if !ok {
		return nil, fmt.Errorf("The snippet %s does not exist", name)
	}
	return execute(partial, name, data)
}

func Subpage(name string, data interface{}) ([]byte, error) {
	subpage, ok := subpages[name]
	if !ok {
		return nil, fmt.Errorf("The subpage %s does not exist", name)
	}
	return execute(subpage, name, data)
}

func execute(t *template.Template, name string, data interface{}) ([]byte, error) {
	var b bytes.Buffer
	t.ExecuteTemplate(&b, name, data)
	return b.Bytes(), nil
}
