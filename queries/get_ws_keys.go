package queries

import (
	"context"

	"github.com/Mahcks/TheGoldenGator/database"
	"go.mongodb.org/mongo-driver/bson"
)

type WebSocketAPIKeys struct {
	User string `json:"user" bson:"user"`
	Key  string `json:"key" bson:"key"`
}

func GetWSKeys() ([]string, error) {
	cursor, err := database.APIKeys.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var apiKeys []WebSocketAPIKeys
	if err = cursor.All(context.Background(), &apiKeys); err != nil {
		return nil, err
	}

	var keys []string
	for _, wsa := range apiKeys {
		keys = append(keys, wsa.Key)
	}

	return keys, nil
}
