package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

type PaginationFooter struct {
	Total   int `json:"total"`
	Current int `json:"current"`
	Results int `json:"results"`
}

func Pagination(limit, page int64) ([]twitch.Member, *PaginationFooter, error) {
	key := "tgg:members"
	limit = limit - 1
	if page <= 0 {
		page = 1
	}

	strMin := strconv.FormatInt((page-1)*limit, 10)
	strLimit := strconv.FormatInt((limit)*page, 10)
	fmt.Println(strMin, strLimit)
	test, err := database.RDB.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: strMin,
		Max: strLimit,
	}).Result()
	if err != nil {
		return nil, nil, err
	}

	var cachedMembers []twitch.Member
	for _, v := range test {
		test2 := twitch.Member{}
		err := json.Unmarshal([]byte(v), &test2)
		if err != nil {
			return nil, nil, err
		}
		cachedMembers = append(cachedMembers, test2)
	}

	totalEntries, err := database.RDB.ZCount(context.Background(), key, "-inf", "+inf").Result()
	if err != nil {
		return nil, nil, err
	}

	totalPages := math.Ceil(float64(totalEntries) / float64(limit))
	footer := PaginationFooter{
		Total:   int(totalPages),
		Current: int(page),
		Results: int(totalEntries),
	}

	return cachedMembers, &footer, nil
}

// Fetches all streamers that are watched for.
func GetStreamers(r *http.Request) ([]twitch.Member, *PaginationFooter, error) {
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		return nil, nil, err
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		return nil, nil, err
	}

	key := "tgg:members"
	cached, err := database.CheckCache(context.Background(), key)
	if err != nil {
		return nil, nil, err
	}

	if cached {
		cachedMembers, footer, err := Pagination(limit, page)
		if err != nil {
			return nil, nil, err
		}
		return cachedMembers, footer, nil
	} else {
		cursor, err := database.Members.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, nil, err
		}

		var users []twitch.Member
		if err = cursor.All(context.Background(), &users); err != nil {
			return nil, nil, err
		}

		// Sort by their login name to make it alphabetical
		sort.Slice(users[:], func(i, j int) bool {
			return users[i].Login < users[j].Login
		})

		// Cache
		for i, m := range users {
			toStore, _ := json.Marshal(m)
			err := database.RDB.ZAdd(context.Background(), key, &redis.Z{
				Score:  float64(i),
				Member: toStore,
			}).Err()

			if err != nil {
				panic(err)
			}
		}

		cachedMembers, footer, err := Pagination(limit, page)
		if err != nil {
			return nil, footer, err
		}

		return cachedMembers, footer, nil
	}
}
