package gochat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000

	TypeUserActivityChannelConnect = "user-channel-connect"
	TypeUserActivityGroupJoin      = "user-group-join"
	TypeUserActivityGroupLeave     = "user-group-leave"
	TypeUserActivityMessageSend    = "user-message-send"
	TypeUserActivityDisconnect     = "user-disconnect"
	TypeUserActivityCustomMessage  = "custom-message"
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

type (
	UserActivity struct {
		Type    string
		User    *User
		Message interface{}
	}

	User struct {
		ID                string            `json:"id" redis:"id"`
		Name              string            `json:"name" redis:"name"`
		AdditionalInfo    map[string]string `json:"additionalInfo,omitempty"`
		onDifferentServer bool
		conn              *websocket.Conn
		server            *Server
		channel           *Channel
		groups            map[string]*Group
		send              chan []byte
		activity          chan *UserActivity
		isActive          bool
		pubSub            *redis.PubSub
	}
)

func NewUser(conn *websocket.Conn, server *Server) *User {
	return &User{
		conn:     conn,
		server:   server,
		groups:   make(map[string]*Group),
		send:     make(chan []byte, 256),
		activity: make(chan *UserActivity),
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
			user.SetActivity(TypeUserActivityChannelConnect, &message)

			user.handleUserConnect(message)
		case CommandMessageSend:
			user.SetActivity(TypeUserActivityMessageSend, &message)

			if user.channel != nil {
				user.handleSendMessage(message)
			}
		case CommandGroupJoin:
			user.SetActivity(TypeUserActivityGroupJoin, &message)

			if user.channel != nil {
				user.handleGroupJoin(message)
			}
		case CommandGroupLeave:
			user.SetActivity(TypeUserActivityGroupLeave, &message)

			if user.channel != nil {
				user.handleGroupLeave(message)
			}
		default:
			user.SetActivity(TypeUserActivityCustomMessage, jsonMessage)
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

func (user *User) GetConn() *websocket.Conn {
	return user.conn
}

func (user *User) Send(message []byte) {
	user.send <- message
}

func (user *User) GetActivity() chan *UserActivity {
	return user.activity
}

func (user *User) SetActivity(activityType string, message interface{}) {
	user.activity <- &UserActivity{
		Type:    activityType,
		User:    user,
		Message: message,
	}
}

func (user *User) ReadActivity() {
	for activity := range user.GetActivity() {
		log.Printf("%s %s", activity.Type, activity.User.ID)
	}
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
				user.server,
			)

			user.server.registerChannel <- user.channel

			go user.channel.Run()
		}

		if !user.isActive {
			user.isActive = true

			go user.subscribePubSub()
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

func (user *User) handleUserdisconnect() {
	user.SetActivity(TypeUserActivityDisconnect, nil)

	if user.channel != nil {
		for _, group := range user.groups {
			group.unregisterUser <- user
		}

		user.channel.unregisterUser <- user
	}

	user.isActive = false

	close(user.send)
	close(user.activity)
	user.conn.Close()

	if user.pubSub != nil {
		user.pubSub.Close()
	}

}

func (user *User) handleGroupJoin(message Message) {
	if message.Group != nil {
		group := user.channel.findGroupByID(message.Group.ID)
		if group == nil {
			group = NewGroup(
				message.Group.ID,
				message.Group.Name,
				message.Group.AdditionalInfo,
			)

			go group.Run()

		}

		user.channel.registerGroup <- group
		group.registerUser <- user
		user.groups[group.ID] = group

		message.User = user
		message.Message = &MessageInfo{
			Type: TypeMessageText,
			Text: MessageGroupJoinSuccessful,
		}
		message.Response = &ResponseInfo{
			Status:  true,
			Message: ResponseMessageSuccess,
		}

		user.send <- []byte(message.encode())
	}
}

func (user *User) handleGroupLeave(message Message) {
	if message.Group != nil {
		group := user.channel.findGroupByID(message.Group.ID)
		if group != nil {
			delete(user.groups, user.ID)
			group.unregisterUser <- user

			if len(group.users) == 0 {
				user.channel.unregisterGroup <- group
			}

			message.User = user
			message.Group = group
			message.Message = &MessageInfo{
				Type: TypeMessageText,
				Text: MessageGroupLeaveSuccessful,
			}
			message.Response = &ResponseInfo{
				Status:  true,
				Message: ResponseMessageSuccess,
			}

			user.send <- []byte(message.encode())
		}

	} else {
		message.Response = &ResponseInfo{
			Status:  false,
			Message: ResponseMessageInvalidPayload,
		}

		user.send <- []byte(message.encode())
	}
}

func (user *User) handleSendMessage(message Message) {
	if message.Message != nil && message.Target != nil {
		switch message.Target.Type {
		case TypeTargetDirect:
			user.handleSendDirectMessage(message)
		case TypeTargetGroup:
			user.handlerSendGroupMessage(message)
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
	if message.Target.User != nil {
		userTarget := user.channel.findUserByID(message.Target.User.ID)
		if userTarget == nil {
			message.Response = &ResponseInfo{
				Status:  false,
				Message: ResponseMessageUserTargetNotConnected,
			}

			message.User = user

			user.send <- []byte(message.encode())
		} else {
			message.User = user
			message.Target.User = userTarget

			if userTarget.onDifferentServer {
				user.publishToPubsub(userTarget.ID, message)
			} else {
				userTarget.send <- []byte(message.encode())
			}

			message.Response = &ResponseInfo{
				Status:  true,
				Message: ResponseMessageSuccess,
			}

			user.send <- []byte(message.encode())
		}
	}
}

func (user *User) handlerSendGroupMessage(message Message) {
	if message.Target.Group != nil {
		groupTarget := user.channel.findGroupByID(message.Target.Group.ID)
		if groupTarget != nil {
			message.User = user
			message.Target.Group = groupTarget

			usersGroupTarget := user.channel.getUsersByGroup(message.Target.Group)
			for _, userGroupTarget := range usersGroupTarget {
				if userGroupTarget.ID != user.ID {
					if userGroupTarget.onDifferentServer {
						user.publishToPubsub(userGroupTarget.ID, message)
					} else {
						userGroupTarget.send <- []byte(message.encode())
					}
				}
			}

			message.Response = &ResponseInfo{
				Status:  true,
				Message: ResponseMessageSuccess,
			}

			user.send <- []byte(message.encode())
		}
	}
}

func (user *User) subscribePubSub() {
	if user.server.PubSub != nil {
		switch user.server.PubSub.driver {
		case PubSubDriverRedis:
			redisClient := user.server.PubSub.conn.(*redis.Client)

			user.pubSub = redisClient.Subscribe(user.server.ctx, fmt.Sprintf("message:%s:%s", user.channel.ID, user.ID))

			for {
				messageBytes, err := user.pubSub.ReceiveMessage(user.server.ctx)
				if err != nil {
					if err.Error() == "redis: client is closed" {
						break
					}

					log.Println(err)
				} else {
					var message Message
					err := json.Unmarshal([]byte(messageBytes.Payload), &message)
					if err != nil {
						log.Println(err)
					} else {
						if user.isActive {
							user.send <- message.encode()
						}
					}
				}
			}
		}
	}
}

func (user *User) publishToPubsub(userTargetID string, message Message) {
	redisClient := user.server.PubSub.conn.(*redis.Client)
	err := redisClient.Publish(user.server.ctx, fmt.Sprintf("message:%s:%s", user.channel.ID, userTargetID), string(message.encode())).Err()
	if err != nil {
		log.Println(err)
	}
}
