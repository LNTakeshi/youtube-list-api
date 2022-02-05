package api

import (
	"net/http"
	"youtubelist/application/service"
)

func (u *Usecase) SpotifyAuth(rw http.ResponseWriter, r *http.Request) {
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	s := service.NewSpotifyService(u.UsecaseBase)
	authURL := s.BeginAuth()
	rw.Header().Set("Location", authURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
