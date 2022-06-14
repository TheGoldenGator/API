package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

const offlineStreams = "tgg:streams:offline"
const onlineStreams = "tgg:streams:online"
const allStreams = "tgg:streams:all"

func StreamPagination(status, sort string, limit, page int64) ([]twitch.PublicStream, *PaginationFooter, error) {
	key := fmt.Sprintf("tgg:streams:%v", strings.ToLower(status))
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

	var cachedStreams []twitch.PublicStream
	for _, v := range test {
		test2 := twitch.PublicStream{}
		err := json.Unmarshal([]byte(v), &test2)
		if err != nil {
			return nil, nil, err
		}
		cachedStreams = append(cachedStreams, test2)
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

	// Sorting by status
	var sortedByStatus []twitch.PublicStream
	if status == "online" || status == "offline" {
		for _, ps := range cachedStreams {
			if ps.Status == status {
				sortedByStatus = append(sortedByStatus, ps)
			}
		}

		sorted := sortStreams(sortedByStatus, sort)
		return sorted, &footer, nil
	} else {
		sorted := sortStreams(cachedStreams, sort)
		return sorted, &footer, nil
	}
}

// Fetches all streams stored in "streams" collection.
func GetStreams(r *http.Request) ([]twitch.PublicStream, *PaginationFooter, error) {
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		return nil, nil, err
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		return nil, nil, err
	}

	// online, offline, all
	status := r.URL.Query().Get("status")

	// viewers, az, za
	sortQ := r.URL.Query().Get("sort")

	key := fmt.Sprintf("tgg:streams:%v", strings.ToLower(status))

	cached, err := database.CheckCache(context.Background(), key)
	if err != nil {
		return nil, nil, err
	}

	if cached {
		cachedStreams, footer, err := StreamPagination(status, sortQ, limit, page)
		if err != nil {
			return nil, nil, err
		}
		return cachedStreams, footer, nil
	} else {
		cursor, err := database.Stream.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, nil, err
		}

		var strims []twitch.PublicStream
		if err = cursor.All(context.Background(), &strims); err != nil {
			return nil, nil, err
		}

		// Sort by their login name to make it alphabetical
		sort.Slice(strims[:], func(i, j int) bool {
			return strims[i].UserLogin < strims[j].UserLogin
		})

		// Cache for online and offline streams
		var onlineCounter int
		var offlineCounter int
		for i := 0; i < len(strims); i++ {
			if strims[i].Status == "online" {
				toStore, _ := json.Marshal(strims[i])
				err := database.RDB.ZAdd(context.Background(), onlineStreams, &redis.Z{
					Score:  float64(onlineCounter),
					Member: toStore,
				}).Err()

				if err != nil {
					panic(err)
				}

				onlineCounter++
			} else if strims[i].Status == "offline" {
				toStore, _ := json.Marshal(strims[i])
				err := database.RDB.ZAdd(context.Background(), offlineStreams, &redis.Z{
					Score:  float64(offlineCounter),
					Member: toStore,
				}).Err()

				if err != nil {
					panic(err)
				}

				offlineCounter++
			}
		}

		for i := 0; i < len(strims); i++ {
			// Add to all streams cache
			toStoreAll, _ := json.Marshal(strims[i])
			err := database.RDB.ZAdd(context.Background(), allStreams, &redis.Z{
				Score:  float64(offlineCounter),
				Member: toStoreAll,
			}).Err()

			if err != nil {
				panic(err)
			}
		}

		cachedStreams, footer, err := StreamPagination(status, sortQ, limit, page)
		if err != nil {
			return nil, footer, err
		}

		return cachedStreams, footer, nil
	}
}

// Clears stream cache to update it.
func GetStreamsDeleteCache() error {
	err := database.RDB.Del(context.Background(), "tgg:streams:online").Err()
	if err != nil {
		return err
	}

	err2 := database.RDB.Del(context.Background(), "tgg:streams:offline").Err()
	if err2 != nil {
		return err2
	}

	err3 := database.RDB.Del(context.Background(), "tgg:streams:all").Err()
	if err3 != nil {
		return err3
	}

	return nil
}

func sortStreams(toParse []twitch.PublicStream, sortMethod string) []twitch.PublicStream {
	if sortMethod == "az" {
		sort.Slice(toParse[:], func(i, j int) bool {
			return toParse[i].UserLogin < toParse[j].UserLogin
		})
	} else if sortMethod == "za" {
		sort.Slice(toParse[:], func(i, j int) bool {
			return toParse[i].UserLogin > toParse[j].UserLogin
		})
	} else if sortMethod == "viewers_high" {
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
	} else if sortMethod == "viewers_low" {
		sort.Slice(toParse, func(i, j int) bool {
			first, _ := strconv.Atoi(toParse[i].StreamViewerCount)
			second, _ := strconv.Atoi(toParse[j].StreamViewerCount)

			if first > second {
				return false
			}
			if first < second {
				return true
			}
			return first < second
		})
	}

	return toParse
}

// Updates stream statuses after every 5 mins
func UpdateStreams() error {
	cursor, err := database.Stream.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	var strims []twitch.PublicStream
	if err = cursor.All(context.Background(), &strims); err != nil {
		return err
	}

	// Sort by their login name to make it alphabetical
	sort.Slice(strims[:], func(i, j int) bool {
		return strims[i].UserLogin < strims[j].UserLogin
	})

	// Cache for online and offline streams
	var onlineCounter int
	var offlineCounter int
	for i := 0; i < len(strims); i++ {
		// Add to all streams cache
		toStoreAll, _ := json.Marshal(strims[i])
		err := database.RDB.ZAdd(context.Background(), allStreams, &redis.Z{
			Score:  float64(offlineCounter),
			Member: toStoreAll,
		}).Err()

		if err != nil {
			panic(err)
		}

		if strims[i].Status == "online" {
			toStore, _ := json.Marshal(strims[i])
			err := database.RDB.ZAdd(context.Background(), onlineStreams, &redis.Z{
				Score:  float64(onlineCounter),
				Member: toStore,
			}).Err()

			if err != nil {
				panic(err)
			}

			onlineCounter++
		} else if strims[i].Status == "offline" {
			toStore, _ := json.Marshal(strims[i])
			err := database.RDB.ZAdd(context.Background(), offlineStreams, &redis.Z{
				Score:  float64(offlineCounter),
				Member: toStore,
			}).Err()

			if err != nil {
				panic(err)
			}

			offlineCounter++
		}
	}

	return nil
}

/*

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

*/
