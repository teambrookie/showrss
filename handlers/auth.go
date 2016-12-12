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

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "empty login/password", http.StatusUnauthorized)
		return
	}
	err, token := betaseries.Auth(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	response := AuthResponse{
		Username: username,
		Token:    token,
	}
	json.NewEncoder(w).Encode(response)
}
