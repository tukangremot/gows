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

	for activity := range user.GetActivity() {
		switch activity.Type {
		case gochat.TypeUserActivityChannelConnect:
			// do somthing when user connect to channel
		case gochat.TypeUserActivityGroupJoin:
			// do somthing when user join to group
		case gochat.TypeUserActivityGroupLeave:
			// do somthing when user leave from group
		case gochat.TypeUserActivityMessageSend:
			// do somthing when user send message
		case gochat.TypeUserActivityDisconnect:
			// do somthing when user disconnet
		}
	}

}

func main() {
	wsServer := gochat.NewServer()
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
