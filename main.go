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
	mux.HandleFunc("/public/", handler.Public)
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/transcribe/", handler.Transcribe)

	fmt.Println("✅ Prattl running")
	fmt.Println(fmt.Sprintf("localhost%s", port))
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
