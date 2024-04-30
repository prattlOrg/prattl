package handler

import (
	"log"
	"net/http"
	"prattl/render"

	"github.com/gorilla/websocket"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templs := [2]string{"index", "recorder"}
	render.RenderTemplate(w, r, templs[:])

}

var upgrader = websocket.Upgrader{}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		t, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %v", string(message))
		// os.WriteFile("content.txt", message, 0666)
		log.Printf("type: %v", t)
		// send base64 encoded string to python
		// audio_bytes = append(audio_bytes, message...)

		// break
	}
}
