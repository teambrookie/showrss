package betaseries

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/teambrookie/showrss/dao"
)

type betaseriesEpisodesResponse struct {
	Shows []struct {
		Unseen []struct {
			ID        int    `json:"id"`
			TheTVDBID int    `json:"thetvdb_id"`
			Title     string `json:"title"`
			Season    int    `json:"season"`
			Episode   int    `json:"episode"`
			Show      struct {
				ID        int    `json:"id"`
				TheTVDBID int    `json:"thetvdb_id"`
				Title     string `json:"title"`
			} `json:"show"`
			Code string `json:"code"`
			User struct {
				Downloaded bool `json:"downloaded"`
			}
		} `json:"unseen"`
	} `json:"shows"`
	Errors []interface{} `json:"errors"`
}

func transformResponse(resp betaseriesEpisodesResponse) []dao.Episode {
	var episodes []dao.Episode
	for _, show := range resp.Shows {
		for _, unseen := range show.Unseen {
			if unseen.User.Downloaded == false {
				episode := dao.Episode{}
				episode.Name = fmt.Sprintf("%s S%02dE%02d", unseen.Show.Title, unseen.Season, unseen.Episode)
				episode.Season = unseen.Season
				episode.Episode = unseen.Episode
				episode.Code = unseen.Code
				episode.ShowID = unseen.Show.TheTVDBID
				episodes = append(episodes, episode)
			}
		}
	}
	return episodes
}

//Episodes retrieve your unseen episode from betaseries
// and flatten the result so you have an array of Episode
func (b Betaseries) Episodes(token string) ([]dao.Episode, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.betaseries.com/episodes/list", nil)
	if err != nil {
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("limit", strconv.Itoa(b.LimitPerShow))
	req.Form = queryParams
	req.Header.Add("X-BetaSeries-Version", "3.0")
	req.Header.Add("X-BetaSeries-Key", b.APIKey)
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
