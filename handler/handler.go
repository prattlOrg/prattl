package handler

import (
	"net/http"
	render "prattl/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "index")
}
