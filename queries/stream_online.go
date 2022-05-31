package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Changes MongoDB status for streamer to online.
func StreamOnline(event twitch.EventSubStreamOnlineEvent) error {
	result, err := database.Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": event.BroadcasterUserID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "status", Value: "online"}}},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Stream went online for %v: %v\n", event.BroadcasterUserLogin, result.ModifiedCount)
	return nil
}
