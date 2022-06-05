package api

import (
	"fmt"
	"net/http"

	"github.com/Mahcks/TheGoldenGator/queries"
	"github.com/Mahcks/TheGoldenGator/websocket"
)

type Message struct {
	// Channel name, or if 'system' if it's a system event
	Channel string `json:"channel"`
	// Event name
	Event string `json:"eventName"`
	// Payload of the event, or if channel is system it'll be an error or message
	Payload string `json:"payload"`
}

func (a *App) serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	key := r.URL.Query().Get("key")
	channels := r.URL.Query().Get("channel")

	if key == "" {
		conn.WriteJSON(Message{
			Channel: "system",
			Event:   "error",
			Payload: "Please provide an API key to connect with",
		})
		conn.Close()
		return
	}

	keys, err := queries.GetWSKeys()
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	isValid := queries.StringInSlice(key, keys)

	if !isValid {
		conn.WriteJSON(Message{
			Channel: "system",
			Event:   "error",
			Payload: "Invalid API key",
		})
		conn.Close()
		return
	}

	if channels == "" {
		conn.WriteJSON(Message{
			Channel: "system",
			Event:   "error",
			Payload: "Please provide a channel to connect to",
		})
		conn.Close()
		return
	} else {
		client := &websocket.Client{
			Channel: channels,
			Conn:    conn,
			Pool:    pool,
		}

		fmt.Println(client.Channel)
		pool.Register <- client
		client.Read()
	}
}
