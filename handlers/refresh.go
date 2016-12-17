package handlers

import (
	"log"
	"net/http"
	"showrss/betaseries"

	"showrss/dao"
)

type refreshHandler struct {
	db   *dao.DB
	jobs chan dao.Episode
}

func (h *refreshHandler) saveAllEpisode(episodes []dao.Episode) error {
	for _, ep := range episodes {
		err := h.db.AddEpisode(ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *refreshHandler) refreshEpisodes(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token must be set in query params", http.StatusNotAcceptable)
		return
	}
	ep, err := betaseries.Episodes(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.saveAllEpisode(ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) refreshTorrent(w http.ResponseWriter, r *http.Request) {
	notFounds, err := h.db.GetAllNotFoundEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, episode := range notFounds {
		h.jobs <- episode
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

func RefreshHandler(db *dao.DB, jobs chan dao.Episode) http.Handler {
	return &refreshHandler{
		db:   db,
		jobs: jobs,
	}
}
