package main

import (
	"flag"
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

func serveWs(server *gochat.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user := gochat.NewUser(conn, server)

	go user.WritePump()
	go user.ReadPump()
}

func main() {
	wsServer := gochat.NewServer(&gochat.Server{
		PubSub: gochat.NewPubSub(),
	})
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
