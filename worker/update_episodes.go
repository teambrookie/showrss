package worker

import (
	"log"
	"strings"

	"github.com/teambrookie/showrss/dao"
)

func UpdateEpisode(datastore *dao.Datastore, episodes <-chan dao.Episode) {
	for episode := range episodes {
		users, err := datastore.GetAllUsers()
		if err != nil {
			log.Println("Error retrieving user from firestore")
		}
		for _, user := range users {
			err := datastore.UpdateUserEpisode(user, episode)
			if err != nil && !strings.Contains(err.Error(), " code = NotFound") {
				log.Println(err)
			}
		}
	}
}
