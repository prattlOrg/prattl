package handler

import (
	"log"
	"net/http"
	"prattl/render"
	"prattl/transcribe"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Public(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/public/", http.FileServer(http.Dir("public/")))
	http.NotFound(w, r)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templs := [2]string{"index", "recorder"}
	render.RenderTemplate(w, r, templs[:])
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		transcribe.TranscribeLocal(string(message))
	}
}
