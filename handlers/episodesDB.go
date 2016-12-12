package handlers

import (
	"encoding/json"
	"net/http"
	"showrss/dao"
)

type dbEpisodeHandler struct {
	db *dao.DB
}

func (h *dbEpisodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	episodes, err := h.db.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func DBEpisodeHandler(db *dao.DB) http.Handler {
	return &dbEpisodeHandler{
		db: db,
	}
}
