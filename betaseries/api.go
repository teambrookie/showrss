package betaseries

import "github.com/fabienfoerster/showrss/dao"

type EpisodeProvider interface {
	Auth(string, string) (string, error)
	Episodes(string) ([]dao.Episode, error)
}

type Betaseries struct {
	ApiKey string
}
