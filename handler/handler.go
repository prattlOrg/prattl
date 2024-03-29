package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"prattl/render"
	"prattl/transcribe"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// render.RenderTemplate(w, r, "index")

	tmpl := template.Must(template.ParseFiles("templates/html/index.html", "templates/html/recorder.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Options(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "options")
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	transcription := transcribe.TranscribeWhisperApi()

	// returning json
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&transcription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
