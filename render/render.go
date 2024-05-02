package render

import (
	"html/template"
	"net/http"
	"path"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, htmlTempls []string) {
	var fpArray []string
	for i := 0; i < len(htmlTempls); i++ {
		fp := path.Join("public/html/", htmlTempls[i]+".html")
		fpArray = append(fpArray, fp)
	}

	tmpl := template.Must(template.ParseFiles(fpArray...))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
