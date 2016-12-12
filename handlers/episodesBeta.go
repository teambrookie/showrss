package handlers

import (
	"encoding/json"
	"net/http"
	"showrss/betaseries"

	"showrss/dao"
)

type betaseriesEpisodeHandler struct {
	db *dao.DB
}

func (h *betaseriesEpisodeHandler) saveAllEpisode(episodes []string) error {
	for _, ep := range episodes {
		episode := dao.Episode{
			Name: ep,
		}
		err := h.db.AddEpisode(episode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *betaseriesEpisodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	err = h.saveAllEpisode(episodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func BetaseriesEpisodeHandler(db *dao.DB) http.Handler {
	return &betaseriesEpisodeHandler{
		db: db,
	}
}
