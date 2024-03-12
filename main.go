package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"prattl/transcribe"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(transcribe.Test()))
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
