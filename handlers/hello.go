package handlers

import "net/http"
import "encoding/json"

type HelloResponse struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Message: "Hello",
	}
	json.NewEncoder(w).Encode(response)
	return
}
