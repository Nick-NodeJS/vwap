package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func Subscribe(conn *websocket.Conn, pair string) error {
	msg := map[string]string{
		MESSAGE_EVENT:   MESSAGE_SUBSCRIBE,
		MESSAGE_CHANNEL: MESSAGE_TRADES,
		MESSAGE_PAIR:    pair,
	}
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}

func ManageSubscriptions(conn *websocket.Conn, pairs []string) {
	for _, pair := range pairs {
		if err := Subscribe(conn, pair); err != nil {
			log.Printf("Failed to subscribe to %s: %v", pair, err)
		}
	}
}
