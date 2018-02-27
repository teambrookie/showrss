package dao

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type Ids struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

type Movie struct {
	Title        string    `json:"title"`
	Year         int       `json:"year"`
	Ids          Ids       `json:"ids"`
	MagnetLink   string    `json:"magnet_link"`
	LastModified time.Time `json:"last_modified"`
	LastAccess   time.Time `json:"last_access"`
}

type MovieStore interface {
	GetMovie(string) (Movie, error)
	AddMovie(Movie) error
	UpdateMovie(Movie) error
	DeleteMovie(string) error
	GetAllMovies() ([]Movie, error)
	GetAllFoundMovies() ([]Movie, error)
	GetAllNotFoundMovies() ([]Movie, error)
}

type BoltMovieStore struct {
	db *bolt.DB
}

func InitDB(dbName string) (*BoltMovieStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltMovieStore{db}, nil
}

func (store *BoltMovieStore) CreateBucket(bucketName string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return err
}

func (store *BoltMovieStore) GetMovie(id string) (Movie, error) {
	var movie Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		v := b.Get([]byte(id))
		json.Unmarshal(v, &movie)
		return nil
	})
	return movie, err

}

func (store *BoltMovieStore) AddMovie(mov Movie) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		if v := b.Get([]byte(strconv.Itoa(mov.Ids.Trakt))); v != nil {
			return nil
		}
		mov.LastAccess = time.Now()
		encoded, err := json.Marshal(mov)
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(mov.Ids.Trakt)), encoded)
	})
	return err
}

func (store *BoltMovieStore) UpdateMovie(mov Movie) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		mov.LastAccess = time.Now()
		encoded, err := json.Marshal(mov)

		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("movies"))
		return b.Put([]byte(strconv.Itoa(mov.Ids.Trakt)), encoded)
	})
	return err
}

func (store *BoltMovieStore) DeleteMovie(id string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		return b.Delete([]byte(id))
	})
	return err
}

func (store *BoltMovieStore) GetAllMovies() ([]Movie, error) {
	var movies []Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		b.ForEach(func(k, v []byte) error {
			var movie Movie
			json.Unmarshal(v, &movie)
			movies = append(movies, movie)
			return nil
		})
		return nil
	})
	return movies, err
}

func (store *BoltMovieStore) GetAllFoundMovies() ([]Movie, error) {
	var movies []Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		b.ForEach(func(k, v []byte) error {
			var movie Movie
			json.Unmarshal(v, &movie)
			if movie.MagnetLink != "" {
				movies = append(movies, movie)
			}

			return nil
		})
		return nil
	})
	return movies, err
}

func (store *BoltMovieStore) GetAllNotFoundMovies() ([]Movie, error) {
	var movies []Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		b.ForEach(func(k, v []byte) error {
			var movie Movie
			json.Unmarshal(v, &movie)
			if movie.MagnetLink == "" {
				movies = append(movies, movie)
			}

			return nil
		})
		return nil
	})
	return movies, err
}
