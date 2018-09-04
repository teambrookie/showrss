package db

import (
	"time"
)

//Media is a generic type
type Media struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Magnet     string      `json:"magnet"`
	Seeds      int         `json:"seeds"`
	Leechs     int         `json:"leechs"`
	LastUpdate time.Time   `json:"last_update"`
	Metadata   interface{} `json:"metadata"`
}

//MediaStore define the interface for retriving media
type MediaStore interface {
	GetCollection(collection string) ([]Media, error)
	GetMedia(mediaID string, collection string) (Media, error)
	AddMedia(media Media, collection string) error
	UpdateMedia(media Media, collection string) error
	DeleteMedia(mediaID string, collection string) error
	Close() error
}
