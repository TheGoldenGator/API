package queries

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Fetches all streamers that are watched for.
func GetStreamers() ([]twitch.Streamer, error) {
	key := "tgg:users"
	cached, err := database.CheckCache(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if cached {
		cached, ok := database.GetCache(context.Background(), key)

		if ok && cached != "" {
			var cS []twitch.Streamer
			json.Unmarshal([]byte(cached), &cS)
			return cS, nil
		}
		return nil, err
	} else {
		cursor, err := database.Users.Find(context.Background(), bson.M{})
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

		errCache := database.SetCache(context.Background(), key, string(toCache), 10*time.Minute)
		if errCache != nil {
			return nil, err
		}

		return users, nil
	}
}
