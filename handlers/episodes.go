package handlers

import (
	"encoding/json"
	"net/http"
	"showrss/dao"
)

type episodeHandler struct {
	db *dao.DB
}

func (h *episodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	episodes, err := h.db.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func EpisodeHandler(db *dao.DB) http.Handler {
	return &episodeHandler{
		db: db,
	}
}
