package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/twitch"
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

				streamerInfo, err := twitch.GetTwitchUser(streamers[i].ID)
				if err != nil {
					return nil, err
				}

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
						TwitchURL:           streamerUrls.TwitchURL,
						RedditURL:           streamerUrls.RedditURL,
						InstagramURL:        streamerUrls.InstagramURL,
						TwitterURL:          streamerUrls.TwitterURL,
						DiscordURL:          streamerUrls.DiscordURL,
						YouTubeURL:          streamerUrls.YouTubeURL,
						TikTokURL:           streamerUrls.TikTokURL,
					}
					insertRes, err := Stream.InsertOne(context.Background(), toInsert)
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
						TwitchURL:           streamerUrls.TwitchURL,
						RedditURL:           streamerUrls.RedditURL,
						InstagramURL:        streamerUrls.InstagramURL,
						TwitterURL:          streamerUrls.TwitterURL,
						DiscordURL:          streamerUrls.DiscordURL,
						YouTubeURL:          streamerUrls.YouTubeURL,
						TikTokURL:           streamerUrls.TikTokURL,
					}
					insertRes, err := Stream.InsertOne(context.Background(), toInsert)
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

func GetStreamerLinks(id string) (twitch.StreamerURLs, error) {
	var search twitch.Streamer
	if err := Users.FindOne(context.Background(), bson.M{"id": id}).Decode(&search); err != nil {
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

// Lowest view count -> highest
/* func lowestViewerCount(streams []twitch.PublicStream) []twitch.PublicStream {
	sort.Slice(streams, func(i, j int) bool {
		if streams[i].StreamViewerCount < streams[j].StreamViewerCount {
			return true
		}
		if streams[i].StreamViewerCount > streams[j].StreamViewerCount {
			return false
		}
		return streams[i].StreamViewerCount < streams[j].StreamViewerCount
	})
} */

// Fetches all streams stored in "streams" collection.
func GetStreams(status, sorted string) ([]twitch.PublicStream, error) {
	var toSearch bson.M
	if status == "online" || status == "offline" {
		toSearch = bson.M{"status": status}
	} else {
		toSearch = bson.M{}
	}

	cursor, err := Stream.Find(context.Background(), toSearch)
	if err != nil {
		return nil, err
	}

	var streams []twitch.PublicStream
	if err = cursor.All(context.Background(), &streams); err != nil {
		return nil, err
	}

	// Sorts based on viewer count
	// Sorts by viewcount: high -> low
	sort.Slice(streams, func(i, j int) bool {
		first, _ := strconv.Atoi(streams[i].StreamViewerCount)
		second, _ := strconv.Atoi(streams[j].StreamViewerCount)

		if first < second {
			return false
		}
		if first > second {
			return true
		}
		return first < second
	})
	return streams, nil
}

// Fetches all streamers that are watched for.
func GetUsers() ([]twitch.Streamer, error) {
	key := "tgg:users"
	cached, err := CheckCache(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if cached {
		cached, ok := GetCache(context.Background(), key)

		if ok && cached != "" {
			var cS []twitch.Streamer
			json.Unmarshal([]byte(cached), &cS)
			return cS, nil
		}
		return nil, err
	} else {
		cursor, err := Users.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, err
		}

		var users []twitch.Streamer
		if err = cursor.All(context.Background(), &users); err != nil {
			return nil, err
		}

		// Cache
		toCache, err := json.Marshal(users)
		if err != nil {
			return nil, err
		}

		errCache := SetCache(context.Background(), key, string(toCache), 10*time.Minute)
		if errCache != nil {
			return nil, err
		}

		return users, nil
	}
}

// Changes MongoDB status for streamer to offline.
func StreamOffline(event twitch.EventSubStreamOfflineEvent) error {
	id, err := strconv.Atoi(event.BroadcasterUserID)
	if err != nil {
		return err
	}

	result, err := Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": id},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "status", Value: "offline"}, {Key: "stream_viewer_count", Value: 0}}},
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
	id, err := strconv.Atoi(event.BroadcasterUserID)
	if err != nil {
		return err
	}

	result, err := Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": id},
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

func ChannelUpdate(event twitch.EventSubChannelUpdateEvent) error {
	id, err := strconv.Atoi(event.BroadcasterUserID)
	if err != nil {
		return err
	}

	result, err := Stream.UpdateOne(
		context.Background(),
		bson.M{"user_id": id},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "stream_title", Value: event.Title}}},
			{Key: "$set", Value: bson.D{{Key: "stream_game_name", Value: event.CategoryName}}},
			{Key: "$set", Value: bson.D{{Key: "stream_game_id", Value: event.CategoryID}}},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Stream changed for %v: %v [%v:%v] changed: %v", event.BroadcasterUserLogin, event.Title, event.CategoryName, event.CategoryID, result.ModifiedCount)
	return nil
}

// This grabs the Twitch team of The Golden Gator(friendzone) and stores them as members
func SortTeamMembers() error {
	fmt.Println("UPDATING TEAM MEMBERS")
	tData, err := twitch.GetTeamMembers()
	if err != nil {
		return err
	}

	// Loops over each member in the team
	t := tData.Data[0]
	for i := 0; i < len(t.Users); i++ {
		// Check if the streamer is in members or not yet.
		var search twitch.Streamer
		if err := Users.FindOne(context.Background(), bson.M{"id": t.Users[i].UserID}).Decode(&search); err != nil {
			if err.Error() == "mongo: no documents in result" {
				// Gets Twitch user data to get the PFP
				twitchUser, err := twitch.GetTwitchUser(t.Users[i].UserID)
				if err != nil {
					return err
				}

				toI := twitch.Streamer{
					ID:              t.Users[i].UserID,
					Login:           t.Users[i].UserLogin,
					DisplayName:     t.Users[i].UserName,
					ProfileImageUrl: twitchUser.Users[0].ProfileImageURL,
					TwitchURL:       fmt.Sprintf("https://www.twitch.tv/%v", t.Users[i].UserLogin),
					InstagramURL:    "N/A",
					RedditURL:       "N/A",
					TwitterURL:      "N/A",
					DiscordURL:      "N/A",
					YouTubeURL:      "N/A",
					TikTokURL:       "N/A",
				}

				Users.InsertOne(context.Background(), toI)
			}
		}
	}

	return nil
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

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

func EventSubscribe() error {
	users, err := GetUsers()
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

func UpdateViewCount() ([]string, error) {
	cursor, err := Stream.Find(context.Background(), bson.M{"status": "online"})
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

		result, err := Stream.UpdateOne(
			context.Background(),
			bson.M{"user_id": streamInfo.Streams[i].ID},
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
