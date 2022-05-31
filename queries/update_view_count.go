package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

/* Polling to update viewer count every 5 minutes */
func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func ViewCountPoll() error {
	_, err := UpdateViewCount()
	if err != nil {
		return err
	}

	return nil
}

func UpdateViewCount() ([]string, error) {
	cursor, err := database.Stream.Find(context.Background(), bson.M{"status": "online"})
	if err != nil {
		return nil, err
	}

	var streams []twitch.PublicStream
	if err = cursor.All(context.Background(), &streams); err != nil {
		return nil, err
	}

	var ids = []string{}
	for i := 0; i < len(streams); i++ {
		ids = append(ids, streams[i].UserID)
	}

	url := "https://api.twitch.tv/helix/streams?first=100&user_id=" + strings.Join(ids, "&user_id=")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+configure.Config.GetString("twitch_client_token"))
	req.Header.Add("Client-Id", configure.Config.GetString("twitch_client_id"))

	c := httpClient()
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	streamInfo := twitch.ManyStreams{}
	if err := json.Unmarshal(body, &streamInfo); err != nil {
		if string(body) == `""` {
			return nil, nil
		}
	}

	for i := 0; i < len(streamInfo.Streams); i++ {
		sViewers := strconv.Itoa(streamInfo.Streams[i].ViewerCount)

		result, err := database.Stream.UpdateOne(
			context.Background(),
			bson.M{"user_id": streamInfo.Streams[i].UserID},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "stream_viewer_count", Value: sViewers}}},
			},
		)

		if err != nil {
			return nil, err
		}

		fmt.Printf("Updated viewer count for %v - %v [%v]\n", streamInfo.Streams[i].UserLogin, streamInfo.Streams[i].ViewerCount, result.ModifiedCount)
	}
	return ids, nil
}
