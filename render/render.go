package render

import (
	"html/template"
	"net/http"
	"path"
)

func RenderTemplate(w http.ResponseWriter, htmlTempls []string) (err error) {
	var fpArray []string
	for i := range htmlTempls {
		fp := path.Join("public/html/", htmlTempls[i]+".html")
		fpArray = append(fpArray, fp)
	}
	tmpl := template.Must(template.ParseFiles(fpArray...))
	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}
