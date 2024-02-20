package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tukangremot/gows"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := gows.NewClient(conn)

	go client.WritePump()
	go client.ReadPump()

	// read message
	for {
		select {
		case message := <-client.ReadMessage():
			// do something whene receive a message
			messageStr := string(message)
			fmt.Println(messageStr)

			// send message to client
			client.SendMessage(message)
		case err := <-client.GetError():
			if err == gows.ErrClientDisconnected {
				// do something if the client disconnects
				fmt.Println(err)
			} else {
				//so something if something else goes wrong
				fmt.Println(err)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", serveWs)

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
