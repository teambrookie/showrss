package worker

import (
	"log"
	"strconv"
	"time"

	"github.com/teambrookie/MediaRSS/movierss/dao"
)

const cleanupInterval = 24 * 1 * time.Hour
const day = 24 * time.Hour
const week = 7 * day

//Cleanup is a function that run periodically and
func Cleanup(out chan<- dao.Movie, store dao.MovieStore) {
	limiter := time.Tick(cleanupInterval)
	for {
		<-limiter
		log.Println("Cleaning movies ...")
		clean(out, store)
	}
}

func clean(out chan<- dao.Movie, store dao.MovieStore) {
	movies, err := store.GetAllFoundMovies()
	if err != nil {
		log.Println("Cleanup routine cannot access movies DB")
	}
	for _, movie := range movies {
		// If you haven't find the movie yet do nothing
		if movie.MagnetLink == "" {
			continue
		}
		// If you have found the movie for more than a week
		if time.Now().Sub(movie.LastModified) > week {
			// and it have been access in the last 12 hours that mean the download never finish so you want a new torrent
			if time.Now().Sub(movie.LastAccess) < day/2 {
				out <- movie
			} else { // if not the download has finished with success and you don't need the movie in the DB anymore ...
				store.DeleteMovie(strconv.Itoa(movie.Ids.Trakt))
			}
		}

	}
}
