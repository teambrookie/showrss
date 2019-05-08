package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fabienfoerster/oauth2client"
	"golang.org/x/oauth2"
)

type oauthHandler struct {
	conf        *oauth2.Config
	newAuthChan chan string
}

func (h *oauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := oauth2client.NewClient(h.conf)
	code := client.RetrieveCode()

	tok, err := h.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// if the auth if successfull , we notify the newAuthChan channel with the user token
	// in a perfect world this would be more difficult but Betaseries token don't expire
	h.newAuthChan <- tok.AccessToken

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tok)

	return

}

func OauthHandler(conf *oauth2.Config, newAuthChan chan string) http.Handler {
	return &oauthHandler{conf, newAuthChan}
}
