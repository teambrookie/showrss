package handlers

import (
	"log"
	"net/http"

	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
)

type refreshHandler struct {
	datastore       *dao.Datastore
	episodeProvider betaseries.EpisodeProvider
	jobs            chan dao.Episode
}

func (h *refreshHandler) refreshEpisodes(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username must be set in query params", http.StatusNotAcceptable)
		return
	}
	token, err := h.datastore.GetUserToken(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	episodes, err := h.episodeProvider.Episodes(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.datastore.AddUserEpisodes(username, episodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO
	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) refreshTorrent(w http.ResponseWriter, r *http.Request) {
	episodes, err := h.datastore.GetAllNotFoundTorrent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, ep := range episodes {
		h.jobs <- ep
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "" && action != "torrent" && action != "episode" {
		http.Error(w, "QueryParam action must be set to torrent or episode", http.StatusMethodNotAllowed)
		return
	}

	if action == "episode" {
		log.Println("Refreshing episodes ...")
		h.refreshEpisodes(w, r)
	}

	if action == "torrent" {
		log.Println("Refreshing torrent ...")
		h.refreshTorrent(w, r)

	}

}

//RefreshHandler handle the refreshing of episodes from Betaseries and the refreshing of torrents using Rarbg
func RefreshHandler(datastore *dao.Datastore, epProvider betaseries.EpisodeProvider, jobs chan dao.Episode) http.Handler {
	return &refreshHandler{
		datastore:       datastore,
		episodeProvider: epProvider,
		jobs:            jobs,
	}
}
