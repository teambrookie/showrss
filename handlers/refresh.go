package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"showrss/betaseries"
	"showrss/torrent"

	"showrss/dao"
)

type refreshHandler struct {
	db *dao.DB
}

func (h *refreshHandler) saveAllEpisode(episodes []string) error {
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
	episodes, err := h.db.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func (h *refreshHandler) refreshTorrent(w http.ResponseWriter, r *http.Request) {
	notFounds, err := h.db.GetAllNotFoundEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, episode := range notFounds {
		torrentLink, err := torrent.Search(episode.Name)
		if err != nil {
			log.Println(err)
		}
		episode.MagnetLink = torrentLink
		h.db.AddEpisode(episode)
	}

	episodes, err := h.db.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
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

func RefreshHandler(db *dao.DB) http.Handler {
	return &refreshHandler{
		db: db,
	}
}
