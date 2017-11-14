package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

type Episode struct {
	Name         string    `json:"name" firestore:"name"`
	Code         string    `json:"code" firestore:"code"`
	ShowID       int       `json:"show_id" firestore:"show_id"`
	MagnetLink   string    `json:"magnet_link" firestore:"magnet_link"`
	LastModified time.Time `json:"last_modified" firestore:"last_modified"`
}

type User struct {
	Username string `firestore:"username"`
	Token    string `firestore:"token"`
}

//Datastore is a struct performing all the shit with the firestore database
type Datastore struct {
	Store *firestore.Client
}

func (d *Datastore) CreateUser(username, token string) error {
	user := User{username, token}
	userRef := d.Store.Collection("users").Doc(username)
	_, err := userRef.Create(context.Background(), user)
	return err
}

func (d *Datastore) GetAllUsers() ([]User, error) {
	var users []User
	docs, err := d.Store.Collection("users").Documents(context.Background()).GetAll()
	fmt.Println(docs)
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var user User
		doc.DataTo(&user)
		users = append(users, user)
	}
	return users, nil
}

func (d *Datastore) existsRef(ref string) bool {
	_, err := d.Store.Doc(ref).Get(context.Background())
	return err == nil
}

func (d *Datastore) AddUserEpisodes(username string, episodes []Episode) error {

	var oldEpisodes = make(map[Episode]bool)
	var newEpisodes = make(map[Episode]bool)

	old, err := d.GetUserEpisodes(username)
	if err != nil {
		return err
	}

	for _, ep := range old {
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

	batch := d.Store.Batch()
	for k, v := range newEpisodes {
		if v {

			episodeRef := d.Store.Collection("users").Doc(username).Collection("episodes").Doc(k.Name)
			batch.Set(episodeRef, k)
			if !d.existsRef("foundTorrents/" + k.Name) {
				newTorrentRef := d.Store.Collection("notFoundTorrents").Doc(k.Name)
				batch.Set(newTorrentRef, k)
			}
		}
	}
	for k, v := range oldEpisodes {
		if v {
			episodeRef := d.Store.Collection("users").Doc(username).Collection("episodes").Doc(k.Name)
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

func (d *Datastore) UpdateUserEpisode(user User, ep Episode) error {
	log.Println("coucou")
	epRef := d.Store.Collection("users").Doc(user.Username).Collection("episodes").Doc(ep.Name)
	res, err := epRef.UpdateStruct(context.Background(), []string{"magnet_link"}, Episode{MagnetLink: ep.MagnetLink})
	log.Println(res)
	return err
}

func (d *Datastore) GetUserEpisodes(username string) ([]Episode, error) {
	var episodes []Episode
	docs, err := d.Store.Collection("users").Doc(username).Collection("episodes").Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var ep Episode
		doc.DataTo(&ep)
		episodes = append(episodes, ep)
	}
	return episodes, nil
}

func (d *Datastore) GetAllNotFoundTorrent() ([]Episode, error) {
	var episodes []Episode
	docs, err := d.Store.Collection("notFoundTorrents").Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var ep Episode
		if err = doc.DataTo(&ep); err != nil {
			return nil, err
		}
		episodes = append(episodes, ep)
	}
	return episodes, nil
}

func (d *Datastore) GetUserToken(username string) (string, error) {
	var user User
	data, err := d.Store.Collection("users").Doc(username).Get(context.Background())
	if err != nil {
		return "", err
	}
	err = data.DataTo(&user)
	if err != nil {
		return "", nil
	}
	return user.Token, nil
}
