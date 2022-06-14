package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStreamerLinks() error {
	cursor, err := database.Stream.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
		return err
	}

	var currStreams []twitch.PublicStream
	if err = cursor.All(context.Background(), &currStreams); err != nil {
		panic(err)
		return err
	}

	for _, stream := range currStreams {
		var streamerBio []twitch.Member
		cursor, err := database.Members.Find(context.Background(), bson.M{"login": stream.UserLogin})

		if err != nil {
			panic(err)
			return err
		}

		if err = cursor.All(context.Background(), &streamerBio); err != nil {
			panic(err)
			return err
		}

		result, err := database.Stream.UpdateOne(
			context.Background(),
			bson.M{"user_login": stream.UserLogin},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "reddit", Value: streamerBio[0].RedditURL},
					{Key: "instagram", Value: streamerBio[0].InstagramURL},
					{Key: "twitter", Value: streamerBio[0].TwitterURL},
					{Key: "discord", Value: streamerBio[0].DiscordURL},
					{Key: "youtube", Value: streamerBio[0].YouTubeURL},
					{Key: "tiktok", Value: streamerBio[0].TikTokURL},
					{Key: "vrchat_legends", Value: streamerBio[0].VRChatLegendsURL},
				}},
			},
		)

		if err != nil {
			panic(err)
			return err
		}
		fmt.Println(result)
	}
	return nil
}
