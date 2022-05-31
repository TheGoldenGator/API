package queries

import (
	"context"
	"sort"
	"strconv"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// Fetches all streams stored in "streams" collection.
func GetStreams(status, sorted string) ([]twitch.PublicStream, error) {
	var toSearch bson.M
	if status == "online" || status == "offline" {
		toSearch = bson.M{"status": status}
	} else {
		toSearch = bson.M{}
	}

	cursor, err := database.Stream.Find(context.Background(), toSearch)
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
