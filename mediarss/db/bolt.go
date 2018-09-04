package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

//BoltMediaStore implement the MediaStore interface using Bolt as the database
type BoltMediaStore struct {
	db *bolt.DB
}

func createBucket(db *bolt.DB, name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
	return err
}

//Open open a new database connection
func Open(dbName string) (*BoltMediaStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err = createBucket(db, "NOTFOUND"); err != nil {
		return nil, err
	}
	if err = createBucket(db, "FOUND"); err != nil {
		return nil, err
	}

	return &BoltMediaStore{db}, nil
}

//Close close the database connection
func (store *BoltMediaStore) Close() error {
	return store.db.Close()
}

//GetCollection retrieve a collection of media
func (store *BoltMediaStore) GetCollection(collection string) ([]Media, error) {
	var medias []Media
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var m Media
			json.Unmarshal(v, &m)
			medias = append(medias, m)
		}
		return nil
	})
	return medias, err
}

//GetMedia retrieve a specified media
func (store *BoltMediaStore) GetMedia(mediaID string, collection string) (Media, error) {
	var media Media
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		v := b.Get([]byte(mediaID))
		json.Unmarshal(v, &media)
		return nil
	})
	return media, err
}

//AddMedia add a media in a specific collection
func (store *BoltMediaStore) AddMedia(media Media, collection string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if v := b.Get([]byte(media.ID)); v != nil {
			return nil
		}
		encoded, err := json.Marshal(media)
		if err != nil {
			return err
		}
		return b.Put([]byte(media.ID), encoded)
	})
	return err
}

//UpdateMedia update a spacific media in a specific collection
func (store *BoltMediaStore) UpdateMedia(media Media, collection string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		encoded, err := json.Marshal(media)
		if err != nil {
			return err
		}
		return b.Put([]byte(media.ID), encoded)
	})
	return err
}

//DeleteMedia delete a media from a specific collection
func (store *BoltMediaStore) DeleteMedia(mediaID string, collection string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		return b.Delete([]byte(mediaID))
	})
	return err
}
