package betaseries

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var apiKey string

type betaseriesAuthResp struct {
	Token string `json:"token"`
}

func toMD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	return checksum
}

//Auth authentificate you in Betaseries using the login/password method
// and return your Betaseries token
func (b Betaseries) Auth(login, password string) (string, error) {
	md5 := toMD5(password)

	data := url.Values{}
	data.Set("login", login)
	data.Add("password", md5)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.betaseries.com/members/auth", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("X-BetaSeries-Version", "2.4")
	req.Header.Add("X-BetaSeries-Key", b.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var betaResp betaseriesAuthResp
	err = json.Unmarshal(body, &betaResp)
	if err != nil {
		return "", err
	}
	return betaResp.Token, nil
}
