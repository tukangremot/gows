package gows

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait          = 10 * time.Second
	pongWait           = 60 * time.Second
	pingPeriod         = (pongWait * 9) / 10
	maxMessageSize     = 10000
	ClientDisconnected = "client is disconnected"
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

type (
	ClientActivity struct {
		Type    string
		Client  *Client
		Message interface{}
	}

	Client struct {
		conn *websocket.Conn
		send chan []byte
		read chan []byte
		err  chan error
	}
)

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
		read: make(chan []byte, 256),
		err:  make(chan error),
	}
}

func (client *Client) ReadPump() {
	defer func() {
		close(client.send)
		client.conn.Close()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			client.err <- err
			break
		}

		client.send <- jsonMessage
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Server closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) GetConn() *websocket.Conn {
	return client.conn
}

func (client *Client) SendMessage(message []byte) {
	client.send <- message
}

func (client *Client) ReadMessage() chan []byte {
	return client.read
}

func (client *Client) GetError() chan error {
	return client.err
}
