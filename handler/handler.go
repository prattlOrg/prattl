package handler

import (
	"fmt"
	"net/http"
	render "prattl/render"
	"prattl/transcribe"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templs := [2]string{"index", "recorder"}
	render.RenderTemplate(w, r, templs[:])
}

// Should accept file form
func Transcribe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("transcribing...")
	transcribe.TranscribeLocal()

	// // getting file from multipart form
	// file, fileHeader, err := r.FormFile("file")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer file.Close()

	// fmt.Println("\nfile:", file, "\nheader:", fileHeader, "\nerr", err)

	// transcription := transcribe.TranscribeWhisperApi(file)

	// // returning json
	// w.Header().Set("Content-Type", "application/json")
	// err = json.NewEncoder(w).Encode(&transcription)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
