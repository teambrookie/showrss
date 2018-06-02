package handlers

import (
	"log"
	"net/http"

	"github.com/teambrookie/MediaRSS/movierss/dao"
	"github.com/teambrookie/MediaRSS/movierss/trakt"
)

type refreshHandler struct {
	store         dao.MovieStore
	movieProvider trakt.MovieProvider
	jobs          chan dao.Movie
}

func (h *refreshHandler) saveAllMovies(movies []dao.Movie) error {
	for _, mov := range movies {
		err := h.store.AddMovie(mov)
		if err != nil {
			log.Println("caca")
			return err
		}
	}
	return nil
}

func (h *refreshHandler) refreshMovies(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "slug must be set in query params", http.StatusNotAcceptable)
		return
	}
	movies, err := h.movieProvider.WatchList(slug, "notCollected")
	err = h.saveAllMovies(movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) refreshTorrent(w http.ResponseWriter, r *http.Request) {
	notFounds, err := h.store.GetAllNotFoundMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, movie := range notFounds {
		h.jobs <- movie
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "" && action != "movie" && action != "torrent" {
		http.Error(w, "QueryParam action must be set to movie or torrent", http.StatusMethodNotAllowed)
		return
	}
	if action == "movie" {
		log.Println("Refreshing movies ...")
		h.refreshMovies(w, r)
	}
	if action == "torrent" {
		log.Println("Refreshing torrent ...")
		h.refreshTorrent(w, r)
	}

}

func RefreshHandler(store dao.MovieStore, movProvider trakt.MovieProvider, jobs chan dao.Movie) http.Handler {
	return &refreshHandler{
		store:         store,
		movieProvider: movProvider,
		jobs:          jobs,
	}
}
