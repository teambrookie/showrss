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
}

type Episode struct {
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	ShowID       int       `json:"show_id"`
	MagnetLink   string    `json:"magnet_link"`
	LastModified time.Time `json:"last_modified"`
}
type Episodes []Episode

type BoltEpisodeStore struct {
	db *bolt.DB
}

func InitDB(dbName string) (*BoltEpisodeStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltEpisodeStore{db}, nil
}

func (store *BoltEpisodeStore) CreateBucket(bucketName string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return err
}

func (store *BoltEpisodeStore) GetEpisode(name string) (Episode, error) {
	var episode Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		v := b.Get([]byte(name))
		json.Unmarshal(v, &episode)
		return nil
	})
	return episode, err

}

func (store *BoltEpisodeStore) AddEpisode(ep Episode) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		if v := b.Get([]byte(ep.Name)); v != nil {
			return nil
		}
		encoded, err := json.Marshal(ep)
		if err != nil {
			return err
		}

		return b.Put([]byte(ep.Name), encoded)
	})
	return err
}

func (store *BoltEpisodeStore) UpdateEpisode(ep Episode) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(ep)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("episodes"))
		return b.Put([]byte(ep.Name), encoded)
	})
	return err
}

func (store *BoltEpisodeStore) DeleteEpisode(name string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		return b.Delete([]byte(name))
	})
	return err
}

func (store *BoltEpisodeStore) GetAllEpisode() ([]Episode, error) {
	var episodes []Episode
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
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
		b := tx.Bucket([]byte("episodes"))
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
