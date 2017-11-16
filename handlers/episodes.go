package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/teambrookie/showrss/dao"
)

type episodeHandler struct {
	datastore *dao.Datastore
}

func (h *episodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]
	episodes, err := h.datastore.GetUserEpisodes(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func EpisodeHandler(datastore *dao.Datastore) http.Handler {
	return &episodeHandler{
		datastore: datastore,
	}
}
