package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mahcks/golden-gator-api/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Creates a stream document for a streamer that doesn't exist
func CreateStream() ([]twitch.Streamer, error) {
	cursor, err := Users.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var streamers []twitch.Streamer
	if err = cursor.All(context.Background(), &streamers); err != nil {
		return nil, err
	}

	for i := 0; i < len(streamers); i++ {
		var str twitch.PublicStream
		if err := Stream.FindOne(context.Background(), bson.M{"user_id": streamers[i].ID}).Decode(&str); err != nil {
			if err.Error() == "mongo: no documents in result" {
				// No document found so create streamer.
				sId := strconv.Itoa(streamers[i].ID)

				streamerInfo, err := twitch.GetTwitchUser(sId)
				if err != nil {
					return nil, err
				}

				uInfo := streamerInfo.Users[0]
				uID, err := strconv.Atoi(uInfo.ID)
				if err != nil {
					return nil, err
				}

				streamData, err := twitch.GetStreamInfo(streamerInfo.Users[0])
				fmt.Println(streamData, err)
				if err != nil {
					return nil, err
				}

				// No streams found with that streamer which means they are offline.
				// No way to get their previous stream data so put "N/A" for things that can't be fetched.
				if len(streamData.Streams) == 0 {
					toInsert := twitch.PublicStream{
						Status:              "offline",
						UserID:              uID,
						UserLogin:           uInfo.Login,
						UserDisplayName:     uInfo.DisplayName,
						UserProfileImageUrl: streamerInfo.Users[0].ProfileImageURL,
						StreamID:            "N/A",
						StreamTitle:         "N/A",
						StreamGameName:      "N/A",
						StreamViewerCount:   0,
						StreamThumbnailUrl:  fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-{width}x{height}.jpg", uInfo.Login),
					}
					insertRes, err := Stream.InsertOne(context.Background(), toInsert)
					if err != nil {
						return nil, err
					}
					fmt.Printf("No stream document found for %s and they are offline so inserting blank document: %s \n", uInfo.Login, insertRes.InsertedID)
					return nil, err
				} else {
					// Stream online and data found, inserting that data.
					toInsert := twitch.PublicStream{
						Status:              "online",
						UserID:              uID,
						UserLogin:           uInfo.Login,
						UserDisplayName:     uInfo.DisplayName,
						UserProfileImageUrl: streamerInfo.Users[0].ProfileImageURL,
						StreamID:            streamData.Streams[0].ID,
						StreamTitle:         streamData.Streams[0].Title,
						StreamGameName:      streamData.Streams[0].GameName,
						StreamViewerCount:   streamData.Streams[0].ViewerCount,
						StreamThumbnailUrl:  fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-{width}x{height}.jpg", uInfo.Login),
					}
					insertRes, err := Stream.InsertOne(context.Background(), toInsert)
					if err != nil {
						return nil, err
					}
					fmt.Println("NO DOCUMENTS FOUND, INSERTED ONE: ", insertRes.InsertedID)
					return nil, err
				}
			}
		}
	}
	return streamers, nil
}

// Fetches all streams stored in "streams" collection.
func GetStreams() ([]twitch.PublicStream, error) {
	cursor, err := Stream.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var streams []twitch.PublicStream
	if err = cursor.All(context.Background(), &streams); err != nil {
		return nil, err
	}

	return streams, nil
}

// Fetches all streamers that are watched for.
func GetUsers() ([]twitch.Streamer, error) {
	cursor, err := Users.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var users []twitch.Streamer
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Changes MongoDB status for streamer to offline.
func StreamOffline(event twitch.EventSubStreamOfflineEvent) error {
	result, err := Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": event.BroadcasterUserID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "status", Value: "offline"}}},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Stream went offline for %v: %v\n", event.BroadcasterUserLogin, result.ModifiedCount)
	return nil
}

// Changes MongoDB status for streamer to online.
func StreamOnline(event twitch.EventSubStreamOnlineEvent) error {
	result, err := Stream.UpdateOne(
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
