package queries

import (
	"context"
	"fmt"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// This grabs the Twitch team of The Golden Gator(friendzone) and stores them as members
func SortTeamMembers() error {
	tData, err := twitch.GetTeamMembers()
	if err != nil {
		return err
	}

	// Loops over each member in the team
	t := tData.Data[0]
	for i := 0; i < len(t.Users); i++ {
		// Check if the streamer is in members or not yet.
		var search twitch.Streamer
		if err := database.Users.FindOne(context.Background(), bson.M{"id": t.Users[i].UserID}).Decode(&search); err != nil {
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

				database.Users.InsertOne(context.Background(), toI)
			}
		}
	}

	return nil
}
