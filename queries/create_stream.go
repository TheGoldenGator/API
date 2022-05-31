package queries

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Creates a stream document for a streamer that doesn't exist
func CreateStream() ([]twitch.Streamer, error) {
	cursor, err := database.Users.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var streamers []twitch.Streamer
	if err = cursor.All(context.Background(), &streamers); err != nil {
		return nil, err
	}

	for i := 0; i < len(streamers); i++ {
		var str twitch.PublicStream
		if err := database.Stream.FindOne(context.Background(), bson.M{"user_id": streamers[i].ID}).Decode(&str); err != nil {
			if err.Error() == "mongo: no documents in result" {
				// No document found so create streamer.

				streamerInfo, err := twitch.GetTwitchUser(streamers[i].ID)
				if err != nil {
					return nil, err
				}

				fmt.Println(streamerInfo)
				uInfo := streamerInfo.Users[0]

				streamData, err := twitch.GetStreamInfo(streamerInfo.Users[0])
				fmt.Println(streamData, err)
				if err != nil {
					return nil, err
				}

				streamerUrls, err := GetStreamerLinks(streamers[i].ID)
				if err != nil {
					return nil, err
				}

				// No streams found with that streamer which means they are offline.
				// No way to get their previous stream data so put "N/A" for things that can't be fetched.
				if len(streamData.Streams) == 0 {
					toInsert := twitch.PublicStream{
						Status:              "offline",
						UserID:              streamers[i].ID,
						UserLogin:           uInfo.Login,
						UserDisplayName:     uInfo.DisplayName,
						UserProfileImageUrl: streamerInfo.Users[0].ProfileImageURL,
						StreamID:            "N/A",
						StreamTitle:         "N/A",
						StreamGameID:        "N/A",
						StreamGameName:      "N/A",
						StreamViewerCount:   "0",
						StreamThumbnailUrl:  fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-{width}x{height}.jpg", uInfo.Login),
						StreamStartedAt:     GetRFCTimestamp(),
						TwitchURL:           streamerUrls.TwitchURL,
						RedditURL:           streamerUrls.RedditURL,
						InstagramURL:        streamerUrls.InstagramURL,
						TwitterURL:          streamerUrls.TwitterURL,
						DiscordURL:          streamerUrls.DiscordURL,
						YouTubeURL:          streamerUrls.YouTubeURL,
						TikTokURL:           streamerUrls.TikTokURL,
					}
					insertRes, err := database.Stream.InsertOne(context.Background(), toInsert)
					if err != nil {
						return nil, err
					}
					fmt.Printf("No stream document found for %s and they are offline so inserting blank document: %s \n", uInfo.Login, insertRes.InsertedID)
				} else {
					// Stream online and data found, inserting that data.
					viewerCountStr := strconv.Itoa(streamData.Streams[0].ViewerCount)
					toInsert := twitch.PublicStream{
						Status:              "online",
						UserID:              streamerInfo.Users[0].ID,
						UserLogin:           uInfo.Login,
						UserDisplayName:     uInfo.DisplayName,
						UserProfileImageUrl: streamerInfo.Users[0].ProfileImageURL,
						StreamID:            streamData.Streams[0].ID,
						StreamTitle:         streamData.Streams[0].Title,
						StreamGameID:        streamData.Streams[0].GameID,
						StreamGameName:      streamData.Streams[0].GameName,
						StreamViewerCount:   viewerCountStr,
						StreamThumbnailUrl:  fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-{width}x{height}.jpg", uInfo.Login),
						StreamStartedAt:     streamData.Streams[0].StartedAt.Format(time.RFC3339),
						TwitchURL:           streamerUrls.TwitchURL,
						RedditURL:           streamerUrls.RedditURL,
						InstagramURL:        streamerUrls.InstagramURL,
						TwitterURL:          streamerUrls.TwitterURL,
						DiscordURL:          streamerUrls.DiscordURL,
						YouTubeURL:          streamerUrls.YouTubeURL,
						TikTokURL:           streamerUrls.TikTokURL,
					}
					insertRes, err := database.Stream.InsertOne(context.Background(), toInsert)
					if err != nil {
						return nil, err
					}
					fmt.Println("NO DOCUMENTS FOUND, INSERTED ONE: ", insertRes.InsertedID)
				}
			}
		}
	}
	return streamers, nil
}
