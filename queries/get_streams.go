package queries

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Fetches all streams stored in "streams" collection.
func GetStreams(status, sorted string) ([]twitch.PublicStream, error) {
	key := "tgg:streams"
	cached, err := database.CheckCache(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if cached {
		cached, ok := database.GetCache(context.Background(), key)
		if ok && cached != "" {
			var data []twitch.PublicStream
			json.Unmarshal([]byte(cached), &data)

			var toParse []twitch.PublicStream
			if status == "online" || status == "offline" {
				for i := 0; i < len(data); i++ {
					if data[i].Status == status {
						toParse = append(toParse, data[i])
					}
				}
			} else {
				data = toParse
			}

			sort := sortHightoLow(toParse, status)
			return sort, nil
		}
		return nil, err
	} else {
		cursor, err := database.Stream.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, err
		}

		var streams []twitch.PublicStream
		if err = cursor.All(context.Background(), &streams); err != nil {
			return nil, err
		}

		toCache, err := json.Marshal(streams)
		if err != nil {
			return nil, err
		}

		errCache := database.RDB.Set(context.Background(), key, string(toCache), time.Minute*10).Err()
		if errCache != nil {
			return nil, errCache
		}

		var toParse []twitch.PublicStream
		if status == "online" || status == "offline" {
			for i := 0; i < len(streams); i++ {
				if streams[i].Status == status {
					toParse = append(toParse, streams[i])
				}
			}
		} else {
			streams = toParse
		}

		sort := sortHightoLow(toParse, status)
		return sort, nil
	}
}

// Clears stream cache to update it.
func GetStreamsDeleteCache() error {
	err := database.RDB.Del(context.Background(), "tgg:streams").Err()
	if err != nil {
		return err
	}
	return nil
}

func sortHightoLow(toParse []twitch.PublicStream, status string) []twitch.PublicStream {
	sort.Slice(toParse, func(i, j int) bool {
		first, _ := strconv.Atoi(toParse[i].StreamViewerCount)
		second, _ := strconv.Atoi(toParse[j].StreamViewerCount)

		if first < second {
			return false
		}
		if first > second {
			return true
		}
		return first < second
	})

	return toParse
}
