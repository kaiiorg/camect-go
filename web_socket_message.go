package camect_go

import "github.com/coder/websocket"

// webSocketMessage is used internally to pass around raw web socket traffic
type webSocketMessage struct {
	Type websocket.MessageType
	Data []byte
	Err  error
}
