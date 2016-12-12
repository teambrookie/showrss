package handlers

import (
	"encoding/json"
	"net/http"
	"showrss/betaseries"
)

type EpisodeResponse struct {
	Episodes []string
}

func EpisodesHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token must be set in query params", http.StatusNotAcceptable)
		return
	}
	episodes, err := betaseries.Episodes(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}
