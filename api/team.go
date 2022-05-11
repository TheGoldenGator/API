package api

import (
	"net/http"

	"github.com/Mahcks/TheGoldenGator/database"
)

func (a *App) TeamData(w http.ResponseWriter, r *http.Request) {
	err := database.SortTeamMembers()
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, err.Error())
	}

	RespondWithJSON(w, r, http.StatusOK, "OK!")
}
