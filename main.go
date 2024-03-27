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

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/options", handler.Options)
	http.HandleFunc("/transcribe", handler.Transcribe)

	fmt.Println("✅ Prattl running")
	fmt.Println(fmt.Sprintf("localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
