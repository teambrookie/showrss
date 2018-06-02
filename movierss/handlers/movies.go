package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teambrookie/MediaRSS/movierss/dao"
)

type movieHandler struct {
	store dao.MovieStore
}

func (h *movieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	movies, err := h.store.GetAllMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	return
}

func MovieHandler(store dao.MovieStore) http.Handler {
	return &movieHandler{
		store: store,
	}
}
