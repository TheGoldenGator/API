package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MDB     *mongo.Client
	RDB     *redis.Client
	Stream  *mongo.Collection
	Members *mongo.Collection
)

func Connect(mongoURI string) error {
	var err error
	// MongoDB Connection
	mdbClientOptions := options.Client().ApplyURI(mongoURI)
	MDB, err = mongo.Connect(context.TODO(), mdbClientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = MDB.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("[INFO] Connected to MongoDB")

	ggdb := MDB.Database("golden_gator")
	Stream = ggdb.Collection("streams")
	Members = ggdb.Collection("members")

	// Redis Connection
	if configure.Config.GetString("environment") == "dev" {
		RDB = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	} else {
		RDB = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "mypassword",
			DB:       0,
		})
	}

	pong, err := RDB.Ping(RDB.Context()).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if pong == "PONG" {
		fmt.Print("[INFO] Connected to Redis\n")
	}

	return err
}
