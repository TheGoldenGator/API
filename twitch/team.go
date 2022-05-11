package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Mahcks/TheGoldenGator/configure"
)

type Team struct {
	Data []TeamData `json:"data"`
}

type TeamMember struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type TeamData struct {
	Users              []TeamMember `json:"users"`
	BackgroundImageURL string       `json:"background_image_url"`
	Banner             string       `json:"banner"`
	CreatedAt          string       `json:"created_at"`
	UpdatedAt          string       `json:"updated_at"`
	Info               string       `json:"info"`
	ThumbnailURL       string       `json:"thumbnail_url"`
	TeamName           string       `json:"team_name"`
	TeamDisplayName    string       `json:"team_display_name"`
	ID                 string       `json:"id"`
}

func GetTeamMembers() (*Team, error) {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/teams?name=friendzone", nil)
	req.Header.Add("Authorization", "Bearer "+configure.Config.GetString("twitch_client_token"))
	req.Header.Add("Client-Id", configure.Config.GetString("twitch_client_id"))

	c := httpClient()
	res, err := c.Do(req)
	if err != nil {
		fmt.Println("Error getting data from Helix API")
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	tData := Team{}
	if err := json.Unmarshal(body, &tData); err != nil {
		return nil, err
	}

	return &tData, nil
}
