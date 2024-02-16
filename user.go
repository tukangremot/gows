package gochat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

type User struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
	conn           *websocket.Conn
	server         *Server
	channel        *Channel
	send           chan []byte
}

func NewUser(conn *websocket.Conn, server *Server) *User {
	return &User{
		conn:   conn,
		server: server,
		send:   make(chan []byte, 256),
	}
}

func (user *User) ReadPump() {
	defer func() {
		user.handleUserdisconnect()
	}()

	user.conn.SetReadLimit(maxMessageSize)
	user.conn.SetReadDeadline(time.Now().Add(pongWait))
	user.conn.SetPongHandler(func(string) error { user.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := user.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		switch message.Command {
		case CommandUserConnect:
			user.handleUserConnect(message)
		case CommandMessageSend:
			if user.channel != nil {
				user.handleSendMessage(message)
			}
		}
	}

}

func (user *User) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		user.conn.Close()
	}()

	for {
		select {
		case message, ok := <-user.send:
			user.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Server closed the channel.
				user.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := user.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(user.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-user.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			user.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := user.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (user *User) handleUserdisconnect() {
	if user.channel != nil {
		user.channel.unregisterUser <- user
	}

	close(user.send)
	user.conn.Close()
}

func (user *User) handleUserConnect(message Message) {
	if message.User != nil && message.Channel != nil {
		user.ID = message.User.ID
		user.Name = message.User.Name
		user.AdditionalInfo = message.User.AdditionalInfo

		user.channel = user.server.findChannelByID(message.Channel.ID)
		if user.channel == nil {
			user.channel = NewChannel(
				message.Channel.ID,
				message.Channel.Name,
				message.Channel.AdditionalInfo,
			)

			user.server.registerChannel <- user.channel

			go user.channel.Run()
		}

		user.channel.registerUser <- user

		message.Message = &MessageInfo{
			Type: TypeMessageText,
			Text: MessageUserConnectSuccessful,
		}

		message.Response = &ResponseInfo{
			Status:  true,
			Message: ResponseMessageSuccess,
		}

	} else {
		message.Response = &ResponseInfo{
			Status:  false,
			Message: ResponseMessageInvalidPayload,
		}
	}

	user.send <- []byte(message.encode())
}

func (user *User) handleSendMessage(message Message) {
	if message.User != nil && message.Message != nil && message.Target != nil {
		switch message.Target.Type {
		case TypeTargetDirect:
			user.handleSendDirectMessage(message)
		}
	} else {
		message.Response = &ResponseInfo{
			Status:  false,
			Message: ResponseMessageInvalidPayload,
		}

		user.send <- []byte(message.encode())
	}
}

func (user *User) handleSendDirectMessage(message Message) {
	userTarget := user.channel.getUserByID(message.Target.User.ID)
	if userTarget == nil {
		message.Response = &ResponseInfo{
			Status:  false,
			Message: ResponseMessageUserTargetNotConnected,
		}

		user.send <- []byte(message.encode())
	} else {
		userTarget.send <- []byte(message.encode())

		message.Response = &ResponseInfo{
			Status:  true,
			Message: ResponseMessageSuccess,
		}

		user.send <- []byte(message.encode())
	}
}
