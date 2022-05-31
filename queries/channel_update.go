package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

func ChannelUpdate(event twitch.EventSubChannelUpdateEvent) error {
	result, err := database.Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": event.BroadcasterUserID},
		bson.M{"$set": bson.M{"stream_title": event.Title, "stream_game_name": event.CategoryName, "stream_game_id": event.CategoryID}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("[CHANNEL.UPDATE] Stream changed for %v: %v [%v:%v] changed: %v", event.BroadcasterUserLogin, event.Title, event.CategoryName, event.CategoryID, result.ModifiedCount)
	return nil
}
