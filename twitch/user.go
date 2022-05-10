package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Mahcks/TheGoldenGator/configure"
)

// HTTP client to make requests
func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func GetTwitchUser(id string) (*ManyUsers, error) {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?id="+id, nil)
	req.Header.Add("Authorization", "Bearer "+configure.Config.GetString("twitch_client_token"))
	req.Header.Add("Client-Id", configure.Config.GetString("twitch_client_id"))

	c := httpClient()
	response, err := c.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return nil, err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	userInfo := ManyUsers{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func GetStreamInfo(user User) (*ManyStreams, error) {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_id="+user.ID, nil)
	req.Header.Add("Authorization", "Bearer "+configure.Config.GetString("twitch_client_token"))
	req.Header.Add("Client-Id", configure.Config.GetString("twitch_client_id"))

	c := httpClient()
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	streamInfo := ManyStreams{}
	if err := json.Unmarshal(body, &streamInfo); err != nil {
		if string(body) == `""` {
			return nil, nil
		}
	}

	return &streamInfo, nil
}
