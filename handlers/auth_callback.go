package handlers

import (
	"net/http"

	"golang.org/x/oauth2"
)

type authCallbackHandler struct {
	conf        oauth2.Config
	newAuthChan chan string
}

func (h *authCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	tok, err := h.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.newAuthChan <- tok.AccessToken

	http.Redirect(w, r, "/rss/"+tok.AccessToken, 301)

	return
}

func AuthCallbackHandler(conf oauth2.Config, newAuthChan chan string) http.Handler {
	return &authCallbackHandler{conf, newAuthChan}
}
