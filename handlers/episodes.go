package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teambrookie/showrss/dao"
)

type episodeHandler struct {
	store dao.EpisodeStore
}

func (h *episodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	episodes, err := h.store.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func EpisodeHandler(store dao.EpisodeStore) http.Handler {
	return &episodeHandler{
		store: store,
	}
}
