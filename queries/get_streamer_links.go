package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

func GetStreamerLinks(id string) (twitch.StreamerURLs, error) {
	var search twitch.Streamer
	if err := database.Users.FindOne(context.Background(), bson.M{"id": id}).Decode(&search); err != nil {
		fmt.Println("ERROR NO LINKS")
		panic(err)
	}

	toSend := twitch.StreamerURLs{
		TwitchURL:    search.TwitchURL,
		RedditURL:    search.RedditURL,
		InstagramURL: search.InstagramURL,
		TwitterURL:   search.TwitterURL,
		DiscordURL:   search.DiscordURL,
		YouTubeURL:   search.YouTubeURL,
		TikTokURL:    search.TikTokURL,
	}

	return toSend, nil
}
