package queries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mahcks/TheGoldenGator/configure"
)

type EventSubCreate struct {
	Type      string                  `json:"type"`
	Version   string                  `json:"version"`
	Condition EventSubCreateCondition `json:"condition"`
	Transport EventSubCreateTransport `json:"transport"`
}

type EventSubCreateCondition struct {
	BroadcasterUserId string `json:"broadcaster_user_id"`
}

type EventSubCreateTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

func EventSubscribe(r *http.Request) error {
	users, _, err := GetStreamers(r)
	if err != nil {
		return err
	}

	/* events := []string{"channel.update", "stream.online", "stream.offline"} */
	c := httpClient()

	for i := 0; i < len(users); i++ {
		toPost := EventSubCreate{
			Type:    "stream.offline",
			Version: "1",
			Condition: EventSubCreateCondition{
				BroadcasterUserId: users[i].ID,
			},
			Transport: EventSubCreateTransport{
				Method:   "webhook",
				Callback: "https://api.thegoldengator.tv/eventsub",
				Secret:   configure.Config.GetString("twitch_eventsub_secret"),
			},
		}

		jsonData, err := json.Marshal(toPost)
		if err != nil {
			return err
		}

		fmt.Println(string(jsonData))
		req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", "Bearer "+configure.Config.GetString("twitch_client_token"))
		req.Header.Add("Client-Id", configure.Config.GetString("twitch_client_id"))
		req.Header.Add("Content-Type", "application/json")

		res, err := c.Do(req)
		if err != nil {
			return err
		}
		fmt.Println(res)
	}
	return nil
}
