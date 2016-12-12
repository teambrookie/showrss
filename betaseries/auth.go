package betaseries

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var apiKey string

type betaseriesAuthResp struct {
	Token string `json:"token"`
}

func init() {
	apiKey = os.Getenv("BETASERIES_KEY")
	if apiKey == "" {
		log.Fatalln("BETASERIES_KEY must be set in env")
	}
}

func toMD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	return checksum
}

func Auth(login, password string) (error, string) {
	md5 := toMD5(password)

	data := url.Values{}
	data.Set("login", login)
	data.Add("password", md5)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.betaseries.com/members/auth", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err, ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("X-BetaSeries-Version", "2.4")
	req.Header.Add("X-BetaSeries-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	var betaResp betaseriesAuthResp
	err = json.Unmarshal(body, &betaResp)
	if err != nil {
		return err, ""
	}
	return nil, betaResp.Token
}
