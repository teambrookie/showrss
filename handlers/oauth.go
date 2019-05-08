package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fabienfoerster/oauth2client"
	"golang.org/x/oauth2"
)

type oauthHandler struct {
	conf *oauth2.Config
}

func (h *oauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := oauth2client.NewClient(h.conf)
	code := client.RetrieveCode()

	tok, err := h.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tok)

	return

}

func OauthHandler(conf *oauth2.Config) http.Handler {
	return &oauthHandler{conf}
}
