package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"prattl/handler"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	port, portOk := os.LookupEnv("PORT")
	if !portOk {
		log.Fatal("Port not defined")
	}

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/transcribe/", handler.Transcribe)
	mux.HandleFunc("/file_test/", handler.FileUpload)

	fmt.Println("✅ Prattl running")
	fmt.Println(fmt.Sprintf("localhost%s", port))
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
