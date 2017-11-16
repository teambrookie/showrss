package handlers

import (
	"net/http"
	"time"
)

type refreshHandler struct {
	limiter chan<- time.Time
}

func (h *refreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.limiter <- time.Now()
	w.WriteHeader(http.StatusOK)
	return

}

//RefreshHandler handle the refreshing of episodes from Betaseries and the refreshing of torrents using Rarbg
func RefreshHandler(limiter chan<- time.Time) http.Handler {
	return &refreshHandler{
		limiter: limiter,
	}
}
