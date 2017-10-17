package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"

	"github.com/teambrookie/showrss/betaseries"
)

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type authHandler struct {
	episodeProvider betaseries.EpisodeProvider
	client          *firestore.Client
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
	response := User{
		Username: username,
		Token:    token,
	}
	userRef := h.client.Collection("users").Doc(username)
	userRef.Create(context.Background(), response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func AuthHandler(client *firestore.Client, episodeProvider betaseries.EpisodeProvider) http.Handler {
	return &authHandler{
		client:          client,
		episodeProvider: episodeProvider,
	}
}
