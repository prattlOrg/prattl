package handler

import (
	"fmt"
	"log"
	"net/http"
	"prattl/render"
	"prattl/transcribe"

	"github.com/gorilla/websocket"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templs := [2]string{"index", "recorder"}
	err := render.RenderTemplate(w, templs[:])
	if err != nil {
		log.Println("template render error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var upgrader = websocket.Upgrader{}

func reader(ws *websocket.Conn) {
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("reader error:", err)
			return
		}
		transcription, err := transcribe.TranscribeLocal(string(message))
		if err != nil {
			log.Println("transcription error:", err)
			return
		}
		fmt.Printf("Result: %s", transcription)
		ws.WriteMessage(messageType, []byte(transcription))
	}
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ws upgrade error:", err)
		return
	}
	defer ws.Close()
	reader(ws)
}
