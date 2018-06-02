package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"github.com/teambrookie/MediaRSS/movierss/dao"
	"github.com/teambrookie/MediaRSS/movierss/trakt"
)

type rssHandler struct {
	store         dao.MovieStore
	movieProvider trakt.MovieProvider
}

func (h *rssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "slug must be set in query params", http.StatusNotAcceptable)
		return
	}
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "MovieRSS by binou",
		Link:        &feeds.Link{Href: "https://github.com/teambrookie/movierss"},
		Description: "A list of torrent for your Track.tv watchlist",
		Author:      &feeds.Author{Name: "Fabien Foerster", Email: "fabienfoerster@gmail.com"},
		Created:     now,
	}
	movies, err := h.movieProvider.WatchList(slug, "notCollected")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, mov := range movies {
		movie, err := h.store.GetMovie(strconv.Itoa(mov.Ids.Trakt))
		h.store.UpdateMovie(movie)
		if movie.MagnetLink == "" || err != nil {
			continue
		}

		description := fmt.Sprintf(`
			<p>Title : %s</p>
			<p>Magnet : <a href="%s">%s</a></p>
			<p>LastModified : %s</p>
			`, movie.Title, movie.MagnetLink, movie.MagnetLink, movie.LastModified)
		item := &feeds.Item{
			Title:       movie.Title,
			Link:        &feeds.Link{Href: movie.MagnetLink},
			Description: description,
			Created:     movie.LastModified,
		}
		feed.Add(item)
	}

	w.Header().Set("Content-Type", "text/xml")
	feed.WriteRss(w)
	return
}

func RSSHandler(store dao.MovieStore, movieProvider trakt.MovieProvider) http.Handler {
	return &rssHandler{
		store:         store,
		movieProvider: movieProvider,
	}
}
