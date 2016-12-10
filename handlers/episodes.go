package handlers

import "net/http"
import "encoding/json"

type EpisodeResponse struct {
	Message string `json:"message"`
}

func EpisodeHandler(w http.ResponseWriter, r *http.Request) {
	response := EpisodeResponse{
		Message: "Hello",
	}
	json.NewEncoder(w).Encode(response)
	return
}
