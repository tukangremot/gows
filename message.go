package gochat

import (
	"encoding/json"
	"log"
)

const (
	CommandUserConnect                    = "user-connect"
	CommandMessageSend                    = "message-send"
	CommandGroupJoin                      = "group-join"
	CommandGroupLeave                     = "group-leave"
	TypeMessageText                       = "text"
	TypeTargetDirect                      = "direct"
	TypeTargetGroup                       = "group"
	MessageUserConnectSuccessful          = "connected successfully"
	MessageGroupJoinSuccessful            = "join group successful"
	MessageGroupLeaveSuccessful           = "leave group successful"
	ResponseMessageSuccess                = "success"
	ResponseMessageInvalidPayload         = "invalid payload"
	ResponseMessageUserTargetNotConnected = "target user is not connected"
)

type (
	Message struct {
		Command  string        `json:"command,omitempty"`
		Channel  *Channel      `json:"channel,omitempty"`
		Group    *Group        `json:"group,omitempty"`
		User     *User         `json:"user,omitempty"`
		Message  *MessageInfo  `json:"message,omitempty"`
		Target   *TargetInfo   `json:"target,omitempty"`
		Response *ResponseInfo `json:"response,omitempty"`
	}

	MessageInfo struct {
		Type           string            `json:"type"`
		Text           string            `json:"text,omitempty"`
		AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
	}

	TargetInfo struct {
		Type  string `json:"type"`
		User  *User  `json:"user,omitempty"`
		Group *Group `json:"group,omitempty"`
	}

	ResponseInfo struct {
		Status  bool   `json:"status"`
		Message string `json:"message,omitempty"`
	}
)

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
