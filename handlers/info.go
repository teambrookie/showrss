package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teambrookie/showrss/dao"
)

type infoHandler struct {
	store dao.EpisodeStore
}

func (h *infoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		http.Error(w, "filename must be set in query params", http.StatusNotAcceptable)
		return
	}
	episodeInfo, err := h.store.GetEpisodeInfo(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodeInfo)
	return

}

func InfoHandler(store dao.EpisodeStore) http.Handler {
	return &infoHandler{store}
}
