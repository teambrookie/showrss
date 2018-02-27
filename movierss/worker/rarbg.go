package worker

import (
	"log"
	"time"

	"github.com/teambrookie/hermes/movierss/dao"
	"github.com/teambrookie/hermes/movierss/torrent"
)

const apiRateLimit = 2 * time.Second

func Rarbg(in <-chan dao.Movie, out chan<- dao.Movie) {
	for movie := range in {
		time.Sleep(apiRateLimit)
		torrentLink, err := torrent.Search(movie.Ids.Imdb)
		if torrentLink == "" {
			log.Printf("%s NOT FOUND", movie.Title)
		}
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", movie.Title, err)
			continue
		}
		movie.MagnetLink = torrentLink
		movie.LastModified = time.Now()
		out <- movie
	}
}
