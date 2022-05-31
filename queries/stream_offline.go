package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Changes MongoDB status for streamer to offline.
func StreamOffline(event twitch.EventSubStreamOfflineEvent) error {
	result, err := database.Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": event.BroadcasterUserID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "status", Value: "offline"}, {Key: "stream_viewer_count", Value: "0"}}},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Stream went offline for %v: %v\n", event.BroadcasterUserLogin, result.ModifiedCount)
	return nil
}
