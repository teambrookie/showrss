package worker

import (
	"log"

	"github.com/teambrookie/hermes/movierss/dao"
)

func DB(in <-chan dao.Movie, store dao.MovieStore) {
	for movie := range in {
		if err := store.UpdateMovie(movie); err != nil {
			log.Printf("Error saving %s to DB ... \n", movie.Title)
		}
	}
}
