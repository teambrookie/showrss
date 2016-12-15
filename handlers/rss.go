package handlers

import (
	"net/http"
	"showrss/dao"
	"time"

	"github.com/gorilla/feeds"
)

type rssHandler struct {
	db *dao.DB
}

func (h *rssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "ShowRSS by binou",
		Link:        &feeds.Link{Href: "https://github.com/TeamBrookie/showrss"},
		Description: "A list of torrent for your unseen Betaseries episodes",
		Author:      &feeds.Author{Name: "Fabien Foerster", Email: "fabienfoerster@gmail.com"},
		Created:     now,
	}
	episodes, err := h.db.GetAllEpisode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, ep := range episodes {
		if ep.MagnetLink == "" {
			continue
		}
		item := &feeds.Item{
			Title:       ep.Name,
			Link:        &feeds.Link{Href: ep.MagnetLink},
			Description: ep.Name,
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
