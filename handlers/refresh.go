package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
)

type refreshHandler struct {
	client          *firestore.Client
	episodeProvider betaseries.EpisodeProvider
	jobs            chan dao.Episode
}

func (h *refreshHandler) saveAllEpisode(episodes []dao.Episode) error {
	batch := h.client.Batch()
	ctx := context.Background()
	notFoundTorrentsRef := h.client.Collection("notFoundTorrents")
	for _, ep := range episodes {
		t := notFoundTorrentsRef.Doc(ep.Name)
		batch.Create(t, ep)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *refreshHandler) getUserEpisodes(username string) []dao.Episode {
	userEpisodes := fmt.Sprintf("users/%s/episodes", username)
	iter := h.client.Collection(userEpisodes).Documents(context.Background())
	var episodes []dao.Episode
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		// if err != nil {
		// 	return nil, err
		// }
		var ep dao.Episode
		doc.DataTo(&ep)
		episodes = append(episodes, ep)
	}
	return episodes
}

func (h *refreshHandler) updateEpisodes(username string) error {
	token := h.getToken(username)
	episodes, _ := h.episodeProvider.Episodes(token)
	var oldEpisodes = make(map[dao.Episode]bool)
	var newEpisodes = make(map[dao.Episode]bool)

	for _, ep := range h.getUserEpisodes(username) {
		oldEpisodes[ep] = true
	}
	for _, ep := range episodes {
		newEpisodes[ep] = true
	}

	for k := range oldEpisodes {
		if _, ok := newEpisodes[k]; ok {
			newEpisodes[k] = false
			oldEpisodes[k] = false
		}
	}

	batch := h.client.Batch()
	for k, v := range newEpisodes {
		if v {
			episodeRef := h.client.Collection("users").Doc(username).Collection("episodes").Doc(k.Name)
			batch.Create(episodeRef, k)
			newTorrentRef := h.client.Collection("notFoundTorrents").Doc(k.Name)
			batch.Set(newTorrentRef, k)

		}
	}
	for k, v := range oldEpisodes {
		if v {
			episodeRef := h.client.Collection("users").Doc(username).Collection("episodes").Doc(k.Name)
			batch.Delete(episodeRef)
		}
	}
	results, err := batch.Commit(context.Background())
	fmt.Println(results)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (h *refreshHandler) getToken(username string) string {
	var user User

	data, _ := h.client.Doc("users/" + username).Get(context.Background())
	data.DataTo(&user)
	return user.Token
}

func (h *refreshHandler) refreshEpisodes(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username must be set in query params", http.StatusNotAcceptable)
		return
	}
	h.updateEpisodes(username)
	w.WriteHeader(http.StatusOK)
	return
}

func (h *refreshHandler) refreshTorrent(w http.ResponseWriter, r *http.Request) {
	iter := h.client.Collection("notFoundTorrents").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var ep dao.Episode
		if err = doc.DataTo(&ep); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
func RefreshHandler(client *firestore.Client, epProvider betaseries.EpisodeProvider, jobs chan dao.Episode) http.Handler {
	return &refreshHandler{
		client:          client,
		episodeProvider: epProvider,
		jobs:            jobs,
	}
}
