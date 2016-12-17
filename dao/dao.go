package dao

import (
	"encoding/json"

	"time"

	"github.com/boltdb/bolt"
)

type Episode struct {
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	ShowID       int       `json:"show_id"`
	MagnetLink   string    `json:"magnet_link"`
	LastModified time.Time `json:"last_modified"`
}
type Episodes []Episode

type DB struct {
	*bolt.DB
}

func InitDB(dbName string) (*DB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) CreateBucket(bucketName string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return err
}

func (db *DB) GetEpisode(name string) (Episode, error) {
	var episode Episode
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		v := b.Get([]byte(name))
		json.Unmarshal(v, &episode)
		return nil
	})
	return episode, err

}

func (db *DB) AddEpisode(ep Episode) error {
	err := db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(ep)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("episodes"))
		return b.Put([]byte(ep.Name), encoded)
	})
	return err
}

func (db *DB) DeleteEpisode(name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		return b.Delete([]byte(name))
	})
	return err
}

func (db *DB) GetAllEpisode() ([]Episode, error) {
	var episodes []Episode
	err := db.View(func(tx *bolt.Tx) error {
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

func (db *DB) GetAllNotFoundEpisode() ([]Episode, error) {
	var episodes []Episode
	err := db.View(func(tx *bolt.Tx) error {
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
