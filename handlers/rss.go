package handlers

import (
	"net/http"
	"showrss/dao"
	"time"

	"showrss/betaseries"

	"github.com/gorilla/feeds"
)

type rssHandler struct {
	db *dao.DB
}

func (h *rssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token must be set in query params", http.StatusNotAcceptable)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "ShowRSS by binou",
		Link:        &feeds.Link{Href: "https://github.com/TeamBrookie/showrss"},
		Description: "A list of torrent for your unseen Betaseries episodes",
		Author:      &feeds.Author{Name: "Fabien Foerster", Email: "fabienfoerster@gmail.com"},
		Created:     now,
	}
	episodes, err := betaseries.Episodes(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, ep := range episodes {
		episode, err := h.db.GetEpisode(ep.Name)
		if episode.MagnetLink == "" || err != nil {
			continue
		}
		item := &feeds.Item{
			Title:       episode.Name,
			Link:        &feeds.Link{Href: episode.MagnetLink},
			Description: episode.Name,
		}
		feed.Add(item)
	}

	w.Header().Set("Content-Type", "text/xml")
	feed.WriteRss(w)
	return

}

func RSSHandler(db *dao.DB) http.Handler {
	return &rssHandler{
		db: db,
	}
}
