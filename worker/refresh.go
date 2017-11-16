package worker

import (
	"log"
	"time"

	"github.com/teambrookie/showrss/betaseries"

	"github.com/teambrookie/showrss/dao"
)

func refreshEpisodes(user dao.User, datastore *dao.Datastore, provider betaseries.EpisodeProvider) error {
	episodes, err := provider.Episodes(user.Token)
	if err != nil {
		return err
	}

	err = datastore.AddUserEpisodes(user.Username, episodes)
	if err != nil {
		return err
	}
	return nil
}

func refreshTorrent(datastore *dao.Datastore, torrentSearch chan dao.Episode) error {
	episodes, err := datastore.GetAllNotFoundTorrent()
	if err != nil {
		return err
	}
	for _, ep := range episodes {
		torrentSearch <- ep
	}
	return nil
}

func Refresh(limiter <-chan time.Time, torrentSearch chan dao.Episode, datastore *dao.Datastore, provider betaseries.EpisodeProvider) {
	for {
		<-limiter
		log.Println("Refreshing episodes ...")
		users, err := datastore.GetAllUsers()
		if err != nil {
			log.Printf("Error fetching users : %s\n", err.Error())
		}
		for _, user := range users {
			err = refreshEpisodes(user, datastore, provider)
			if err != nil {
				log.Printf("Error refresh episodes for %s : %s\n", user.Username, err.Error())
			}
		}
		err = refreshTorrent(datastore, torrentSearch)
		if err != nil {
			log.Printf("Error refreshing torrents : %s\n", err.Error())
		}

	}
}
