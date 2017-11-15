package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
)

type authHandler struct {
	episodeProvider betaseries.EpisodeProvider
	datastore       *dao.Datastore
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "empty login/password", http.StatusUnauthorized)
		return
	}
	token, err := h.episodeProvider.Auth(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	response := dao.User{
		Username: username,
		Token:    token,
	}
	err = h.datastore.CreateUser(username, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func AuthHandler(datastore *dao.Datastore, episodeProvider betaseries.EpisodeProvider) http.Handler {
	return &authHandler{
		datastore:       datastore,
		episodeProvider: episodeProvider,
	}
}
