package gows

import "errors"

var (
	ErrClientDisconnected = errors.New("client is disconnected")
)

var websocketError = map[string]error{
	"websocket: close 1005 (no status)": ErrClientDisconnected,
	"websocket: close 1000 (normal)":    ErrClientDisconnected,
}

func parseError(err error) error {
	if websocketError[err.Error()] != nil {
		return websocketError[err.Error()]
	}

	return err
}
