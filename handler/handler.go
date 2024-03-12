package handler

import (
	"net/http"
	"prattl/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "index")
}
