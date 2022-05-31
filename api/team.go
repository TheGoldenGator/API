package api

import (
	"net/http"

	"github.com/Mahcks/TheGoldenGator/queries"
)

func (a *App) TeamData(w http.ResponseWriter, r *http.Request) {
	err := queries.SortTeamMembers()
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, err.Error())
	}

	RespondWithJSON(w, r, http.StatusOK, "OK!")
}
