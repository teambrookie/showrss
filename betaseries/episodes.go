package betaseries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	apiKey = os.Getenv("BETASERIES_KEY")
	if apiKey == "" {
		log.Fatalln("BETASERIES_KEY must be set in env")
	}
}

type betaseriesEpisodesResponse struct {
	Shows []struct {
		Unseen []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Season  int    `json:"season"`
			Episode int    `json:"episode"`
			Show    struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
			} `json:"show"`
			Code string `json:"code"`
		} `json:"unseen"`
	} `json:"shows"`
	Errors []interface{} `json:"errors"`
}

func transformResponse(resp betaseriesEpisodesResponse) []string {
	var episodes []string
	for _, show := range resp.Shows {
		for _, unseen := range show.Unseen {
			episode := fmt.Sprintf("%s S%02dE%02d", unseen.Show.Title, unseen.Season, unseen.Episode)
			episodes = append(episodes, episode)
		}

	}
	return episodes
}

//Episodes retrieve your unseen episode from betaseries
func Episodes(token string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.betaseries.com/episodes/list", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-BetaSeries-Version", "2.4")
	req.Header.Add("X-BetaSeries-Key", apiKey)
	req.Header.Add("X-BetaSeries-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var betaResp betaseriesEpisodesResponse
	err = json.NewDecoder(resp.Body).Decode(&betaResp)
	if err != nil {
		return nil, err
	}
	return transformResponse(betaResp), nil
}
