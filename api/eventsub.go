package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/queries"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"github.com/Mahcks/TheGoldenGator/websocket"
)

// Verify message from EventSub
func VerifyEventSubNotification(secret string, header http.Header, message string) bool {
	hmacMessage := []byte(fmt.Sprintf("%s%s%s", header.Get("Twitch-Eventsub-Message-Id"), header.Get("Twitch-Eventsub-Message-Timestamp"), message))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(hmacMessage)
	hmacsha256 := fmt.Sprintf("sha256=%s", hex.EncodeToString(mac.Sum(nil)))
	return hmacsha256 == header.Get("Twitch-Eventsub-Message-Signature")
}

type eventsubNotification struct {
	Subscription twitch.EventSubSubscription `json:"subscription"`
	Challenge    string                      `json:"challenge"`
	Event        json.RawMessage             `json:"event"`
}

// Route that fetches POSTed eventsub notifications from Twitch
func (a *App) EventsubRecievedNotification(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	defer r.Body.Close()

	// Verify Twitch sent the event
	if !VerifyEventSubNotification(configure.Config.GetString("twitch_eventsub_secret"), r.Header, string(body)) {
		log.Println("No valid signature on subscription")
		return
	} else {
		log.Println("Verified signature on subscription")
	}
	var vals eventsubNotification
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if err != nil {
		log.Println(err)
		return
	}

	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))
		return
	}

	eventType := bytes.NewBuffer([]byte(vals.Subscription.Type)).String()

	switch {
	case eventType == "stream.online":
		var streamOnline twitch.EventSubStreamOnlineEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOnline)
		if err != nil {
			panic(err.Error())
		}

		errDb := queries.StreamOnline(streamOnline)
		if errDb != nil {
			panic(err.Error())
		}

	case eventType == "stream.offline":
		var streamOffline twitch.EventSubStreamOfflineEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOffline)
		if err != nil {
			panic(err.Error())
		}

		errDb := queries.StreamOffline(streamOffline)
		if errDb != nil {
			panic(err.Error())
		}

	case eventType == "channel.update":
		var streamUpdate twitch.EventSubChannelUpdateEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamUpdate)
		if err != nil {
			panic(err.Error())
		}

		errDb := queries.ChannelUpdate(streamUpdate)
		if errDb != nil {
			panic(err.Error())
		}

	case eventType == "channel.channel_points_custom_reward_redemption.add":
		var redemption twitch.EventSubChannelPointRewardRedemption
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&redemption)
		if err != nil {
			panic(err.Error())
		}

		payload, err := json.Marshal(redemption)
		if err != nil {
			return
		}
		websocket.PublishEvent(strings.ToLower(redemption.BroadcasterUserLogin), "channel.channel_points_custom_reward_redemption.add", string(payload))
	}
}
