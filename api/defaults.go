package api

import (
	"net/http"

	"github.com/Mahcks/TheGoldenGator/database"
)

type Url struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Default home route
// URLs godoc
// @Summary Returns list of avilable endpoints along with a description.
// @Tags Default
// @Produce json
// @Success 200 {array} Url
// @Router / [get]
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	urls :=
		[]Url{
			{
				Description: "Get all the current Golden Gator streams that are offline and online.",
				URL:         "localhost:8000/streams",
			},
			{
				Description: "Gets all streamers that are listed in the Golden Gator",
				URL:         "localhost:8000/streamers",
			},
		}

	RespondWithJSON(w, r, http.StatusOK, urls)
}

// Fetches all streams that are stored.
// URLs godoc
// @Summary Returns list of all Golden Gator streamers that are online and offline.
// @Description Using this endpoint, you'll be able to get all stored data about their stream and streamer.
// @Tags Default
// @Produce json
// @Success 200 {array} twitch.PublicStream
// @Router /streams [get]
func (a *App) Streams(w http.ResponseWriter, r *http.Request) {
	streams, err := database.GetStreams()
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, err.Error())
	}

	RespondWithJSON(w, r, http.StatusOK, streams)
}

// Fetches all streamers
// URLs godoc
// @Summary Returns list of all Golden Gator streamers.
// @Description Using this endpoint, you'll be able to get all stored data about a streamer.
// @Tags Default
// @Produce json
// @Success 200 {array} twitch.Streamer
// @Router /users [get]
func (a *App) Users(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, err.Error())
	}

	RespondWithJSON(w, r, http.StatusOK, users)
}

// If route doesn't exist then this will be called.
func (a *App) NotFound(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, r, http.StatusNotFound, "This doesn't exist Despair")
}

func (a *App) Test(w http.ResponseWriter, r *http.Request) {
	data, err := database.CreateStream()
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, err.Error())
	}

	RespondWithJSON(w, r, http.StatusOK, data)
}
