package gochat

import (
	"encoding/json"
	"log"
)

const (
	CommandUserConnect                    = "user-connect"
	CommandMessageSend                    = "message-send"
	TypeMessageText                       = "text"
	TypeTargetDirect                      = "direct"
	MessageUserConnectSuccessful          = "connected successfully"
	ResponseMessageSuccess                = "success"
	ResponseMessageUserTargetNotConnected = "target user is not connected"
	ResponseMessagMessageSendSuccessfull  = "send message successful"
	ResponseMessageInvalidPayload         = "invalid payload"
)

type (
	Message struct {
		Command  string        `json:"command,omitempty"`
		Channel  *Channel      `json:"channel,omitempty"`
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
		Type string `json:"type"`
		User *User  `json:"user,omitempty"`
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
