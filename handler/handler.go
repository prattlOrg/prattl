package handler

import (
	"fmt"
	"log"
	"net/http"
	"prattl/render"
	"prattl/transcribe"
	"time"

	"github.com/gorilla/websocket"
)

//// Logging

func trace(s string) (string, time.Time) {
	fmt.Println()
	log.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("END:", s, "took", endTime.Sub(startTime))
	fmt.Println()
}

//// Logging

var upgrader = websocket.Upgrader{}

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

		defer un(trace("transcribe.TranscribeLocal()"))
		fmt.Println("transcribing...")
		transcribe.TranscribeLocal(string(message))
	}
}
