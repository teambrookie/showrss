package handlers

import (
	"net/http"

	"golang.org/x/oauth2"
)

type oauthHandler struct {
	conf *oauth2.Config
}

func (h *oauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	url := h.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, 301)

}

func OauthHandler(conf *oauth2.Config) http.Handler {
	return &oauthHandler{conf}
}
