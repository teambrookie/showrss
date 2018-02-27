package betaseries

import "github.com/teambrookie/hermes/showrss/dao"

// EpisodeProvider is a generic interface for fetching unseen episodes
type EpisodeProvider interface {
	Auth(string, string) (string, error)
	Episodes(string) ([]dao.Episode, error)
}

// Betaseries is a struct that will implement the EpisodeProvider interface
type Betaseries struct {
	APIKey string
}
