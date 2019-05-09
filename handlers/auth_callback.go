package handlers

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type authCallbackHandler struct {
	conf        *oauth2.Config
	newAuthChan chan string
	host        string
}

func (h *authCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	log.Printf("The fucking code is : %s", code)
	tok, err := h.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.newAuthChan <- tok.AccessToken

	redirectURL := fmt.Sprintf("%s/rss/%s", h.host, tok.AccessToken)
	http.Redirect(w, r, redirectURL, 301)

	return
}

func AuthCallbackHandler(conf *oauth2.Config, newAuthChan chan string, host string) http.Handler {
	return &authCallbackHandler{conf, newAuthChan, host}
}
