package dao

import (
	"encoding/json"

	"time"

	"github.com/boltdb/bolt"
)

type EpisodeStore interface {
	GetEpisode(string) (Episode, error)
	AddEpisode(Episode) error
	UpdateEpisode(Episode) error
	DeleteEpisode(string) error
	GetAllEpisode() ([]Episode, error)
	GetAllNotFoundEpisode() ([]Episode, error)
	GetEpisodeInfo(string) (Episode, error)
}

type Episode struct {
	Name         string    `json:"name"`
	Season       int       `json:"season"`
	Episode      int       `json:"episode"`
	Code         string    `json:"code"`
	ShowID       int       `json:"show_id"`
	MagnetLink   string    `json:"magnet_link"`
	TorrentURL   string    `json:"torrent_url"`
	Filename     string    `json:"filename"`
	LastModified time.Time `json:"last_modified"`
}
type Episodes []Episode

const episodesBucket string = "episodes"
const secondaryKeyBucket string = "secondaryKey"

type BoltEpisodeStore struct {
	db *bolt.DB
}

func InitDB(dbName string) (*BoltEpisodeStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(episodesBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(secondaryKeyBucket))
		return err
	})
	return &BoltEpisodeStore{db}, err
}

func (store *BoltEpisodeStore) GetEpisode(name string) (Episode, error) {
	var episode Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucket))
		v := b.Get([]byte(name))
		json.Unmarshal(v, &episode)
		return nil
	})
	return episode, err
}

func (store *BoltEpisodeStore) AddEpisode(ep Episode) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		epBucket := tx.Bucket([]byte(episodesBucket))
		if v := epBucket.Get([]byte(ep.Name)); v != nil {
			return nil
		}
		encoded, err := json.Marshal(ep)
		if err != nil {
			return err
		}
		return epBucket.Put([]byte(ep.Name), encoded)
	})
	return err
}

func (store *BoltEpisodeStore) UpdateEpisode(ep Episode) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(ep)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte(episodesBucket))
		err = b.Put([]byte(ep.Name), encoded)
		if err != nil {
			return err
		}
		secKeyBucket := tx.Bucket([]byte(secondaryKeyBucket))
		return secKeyBucket.Put([]byte(ep.Filename), []byte(ep.Name))
	})
	return err
}

func (store *BoltEpisodeStore) DeleteEpisode(name string) error {
	ep, err := store.GetEpisode(name)
	if err != nil {
		return err
	}
	err = store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucket))
		err := b.Delete([]byte(name))
		if err != nil {
			return err
		}
		secKeyBucket := tx.Bucket([]byte(secondaryKeyBucket))
		return secKeyBucket.Delete([]byte(ep.Filename))
	})
	return err
}

func (store *BoltEpisodeStore) GetEpisodeInfo(filename string) (Episode, error) {

	var episode Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		secKeyBucket := tx.Bucket([]byte(secondaryKeyBucket))
		primaryKey := secKeyBucket.Get([]byte(filename))
		b := tx.Bucket([]byte(episodesBucket))
		v := b.Get([]byte(primaryKey))
		json.Unmarshal(v, &episode)
		return nil
	})
	return episode, err
}

func (store *BoltEpisodeStore) GetAllEpisode() ([]Episode, error) {
	var episodes []Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var episode Episode
			json.Unmarshal(v, &episode)
			episodes = append(episodes, episode)
		}
		return nil
	})
	return episodes, err
}

func (store *BoltEpisodeStore) GetAllNotFoundEpisode() ([]Episode, error) {
	var episodes []Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var episode Episode
			json.Unmarshal(v, &episode)
			if episode.MagnetLink == "" {
				episodes = append(episodes, episode)
			}

		}
		return nil
	})
	return episodes, err
}
