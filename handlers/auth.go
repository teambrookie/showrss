package handlers

import (
	"encoding/json"
	"net/http"
	"showrss/betaseries"
)

type AuthResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type authHandler struct {
	episodeProvider betaseries.EpisodeProvider
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
	response := AuthResponse{
		Username: username,
		Token:    token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func AuthHandler(episodeProvider betaseries.EpisodeProvider) http.Handler {
	return &authHandler{
		episodeProvider: episodeProvider,
	}
}
