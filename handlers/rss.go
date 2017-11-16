package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"fmt"

	"github.com/gorilla/feeds"
	"github.com/teambrookie/showrss/dao"
)

type rssHandler struct {
	datastore *dao.Datastore
}

func (h *rssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]
	if username == "" {
		http.Error(w, "username must be set in query params", http.StatusNotAcceptable)
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
	episodes, err := h.datastore.GetUserEpisodes(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, episode := range episodes {
		if episode.MagnetLink == "" || err != nil {
			continue
		}
		description := fmt.Sprintf("<a href='%s'>MagnetLink</a>", episode.MagnetLink)
		item := &feeds.Item{
			Title:       episode.Name,
			Link:        &feeds.Link{Href: episode.MagnetLink},
			Description: description,
			Created:     episode.LastModified,
		}
		feed.Add(item)
	}

	w.Header().Set("Content-Type", "text/xml")
	feed.WriteRss(w)
	return

}

func RSSHandler(datastore *dao.Datastore) http.Handler {
	return &rssHandler{
		datastore: datastore,
	}
}
