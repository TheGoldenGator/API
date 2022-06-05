package websocket

import (
	"fmt"
	"strings"
)

var WSPool *Pool

func init() {
	WSPool = NewPool()
}

func PublishEvent(channel, event, payload string) {
	for c := range WSPool.Clients {
		channels := strings.Split(c.Channel, ",")

		for i := 0; i < len(channels); i++ {
			if channel == channels[i] {
				fmt.Println(Message{Channel: channel, Event: event, Payload: payload})
				c.Conn.WriteJSON(Message{Channel: channel, Event: event, Payload: payload})
			}
		}
	}
}
