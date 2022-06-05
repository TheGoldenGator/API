package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string
	Channel string
	Conn    *websocket.Conn
	Pool    *Pool
}

type Message struct {
	// Channel name, or if 'system' if it's a system event
	Channel string `json:"channel"`
	// Event name
	Event string `json:"eventName"`
	// Payload of the event, or if channel is system it'll be an error or message
	Payload string `json:"payload"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(p)
		/* message := Message{Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message) */
	}
}
