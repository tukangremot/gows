package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/tukangremot/gows"
)

type User struct {
	ID     string
	client *gows.Client
}

func (u *User) Register(ctx context.Context, rdb *redis.Client) error {
	return rdb.Set(ctx, fmt.Sprintf("u:%s", u.ID), true, 0).Err()
}

func (u *User) Unregister(ctx context.Context, rdb *redis.Client) error {
	return rdb.Del(ctx, fmt.Sprintf("u:%s", u.ID)).Err()
}

func (u *User) BroadcastMessage(ctx context.Context, rdb *redis.Client, message []byte) error {
	keys, err := rdb.Keys(ctx, "u:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		userID := strings.ReplaceAll(key, "u:", "")

		err := rdb.Publish(ctx, fmt.Sprintf("m:%s", userID), message).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) SubscribePubsub(ctx context.Context, rdb *redis.Client) {
	pubsub := rdb.Subscribe(ctx, fmt.Sprintf("m:%s", u.ID))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		u.client.SendMessage([]byte(msg.Payload))
	}
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	addr  = flag.String("addr", ":8080", "http service address")
	users = make(map[string]User)
	rdb   = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := gows.NewClient(conn)

	go client.WritePump()
	go client.ReadPump()

	ctx := r.Context()

	user := User{
		ID:     r.Header.Get("X-User-ID"),
		client: client,
	}

	// register to pubsub
	user.Register(ctx, rdb)

	// register to memory
	users[user.ID] = user

	// user subscribe to pubsub
	go user.SubscribePubsub(ctx, rdb)

	// read message
	for {
		select {
		case message := <-client.ReadMessage():
			user.BroadcastMessage(ctx, rdb, message)
		case err := <-client.GetError():
			if err == gows.ErrClientDisconnected {
				user.Unregister(ctx, rdb)
				delete(users, user.ID)
			} else {
				log.Println(err)
			}
		}
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ok")
	})

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
