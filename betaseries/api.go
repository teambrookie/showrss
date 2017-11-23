package betaseries

import "github.com/teambrookie/showrss/dao"
import "encoding/json"
import "errors"
import "fmt"

// EpisodeProvider is a generic interface for fetching unseen episodes
type EpisodeProvider interface {
	Auth(string, string) (string, error)
	Episodes(string) ([]dao.Episode, error)
}

// Betaseries is a struct that will implement the EpisodeProvider interface
type Betaseries struct {
	APIKey string
}

//Error wrap the error response from the Betaseries API
type betaError struct {
	Errors []struct {
		Code int    `json:"code"`
		Text string `json:"text"`
	} `json:"errors"`
}

//CatchAPIError parse the JSON send by the Betaseries API and check for error message
func CatchAPIError(data []byte) error {
	var apiError betaError
	json.Unmarshal(data, &apiError)
	if apiError.Errors[0].Code != 0 {
		errorText := fmt.Sprintf("Betaseries API Error #%d : %s", apiError.Errors[0].Code, apiError.Errors[0].Text)
		return errors.New(errorText)
	}
	return nil
}
