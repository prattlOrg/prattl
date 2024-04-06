package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	r.ParseMultipartForm(200)
	mForm := r.MultipartForm

	// for k, _ := range mForm.File {
	// 	file, _, _ := r.FormFile(k)
	// 	fmt.Printf("BODY: %s", file.Close)
	// }
	transcription := transcribe.TranscribeWhisperApi(mForm)
	// should pass file form data to below

	// returning json
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&transcription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// File blob testing
func FileUpload(w http.ResponseWriter, r *http.Request) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	r.ParseMultipartForm(32 << 20)
	//ParseMultipartForm parses a request body as multipart/form-data
	file, handler, err := r.FormFile("file") //retrieve the file from form data
	//replace file with the key your sent your image with
	if err != nil {
		return
	}
	defer file.Close() //close the file when we finish
	//this is path which  we want to store the file
	f, err := os.OpenFile("path/to/save/image/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Printf("FILENAME %s", handler.Filename)
	//here we save our file to our path
}
