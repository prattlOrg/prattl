package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"prattle/endpoint"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	fmt.Println(endpoint.Hey())

	log.Fatal(http.ListenAndServe(":8081", nil))
}
