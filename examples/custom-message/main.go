package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tukangremot/gochat"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
payload:

	{
	    "text": "<your-text>"
	}
*/
type Message struct {
	Text string `json:"text"`
}

func serveWs(server *gochat.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user := gochat.NewUser(conn, server)

	go user.WritePump()
	go user.ReadPump()

	for activity := range user.GetActivity() {
		switch activity.Type {
		case gochat.TypeUserActivityDisconnect:
			// do somthing when user disconnet
		default:
			var message *Message

			msgString := activity.Message.([]byte)
			err := json.Unmarshal(msgString, &message)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(message.Text)

			user.Send(msgString) // example: send message to user
		}

	}

}

func main() {
	wsServer := gochat.NewServer(nil)
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsServer, w, r)
	})

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
