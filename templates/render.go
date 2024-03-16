package render

import (
	"html/template"
	"net/http"
	"path"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, html string) {
	fp := path.Join("templates/", html+".html")
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
