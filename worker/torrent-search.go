package worker

import (
	"context"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/teambrookie/showrss/dao"
	"github.com/teambrookie/showrss/torrent"
)

func TorrentSearch(jobs <-chan dao.Episode, updateEpisode chan<- dao.Episode, client *firestore.Client) {
	time.Sleep(2 * time.Second)
	for episode := range jobs {
		time.Sleep(2 * time.Second)
		torrentLink, err := torrent.Search(strconv.Itoa(episode.ShowID), episode.Code, "720p")
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", episode.Name, err)
			continue
		}
		if torrentLink == "" {
			log.Printf("Cannot find : %s", episode.Name)
			continue
		}
		episode.MagnetLink = torrentLink
		episode.LastModified = time.Now()
		batch := client.Batch()
		oldRef := client.Collection("notFoundTorrents").Doc(episode.Name)
		newRef := client.Collection("foundTorrents").Doc(episode.Name)
		batch.Set(newRef, episode)
		batch.Delete(oldRef)
		_, err = batch.Commit(context.Background())
		updateEpisode <- episode
		if err != nil {
			log.Printf("Error saving %s to DB ...\n", episode.Name)
		}

	}
}
