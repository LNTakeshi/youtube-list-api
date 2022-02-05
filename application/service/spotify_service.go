package service

import (
	"context"
	"youtubelist/application/spotify"
	"youtubelist/application/usecase/util"
	"youtubelist/util/log"

	"golang.org/x/oauth2"
)

type SpotifyService struct {
	log     log.Logger
	spotify spotify.ISpotify
}

func NewSpotifyService(base *util.UsecaseBase) *SpotifyService {
	return &SpotifyService{log: base.Log, spotify: base.Spotify}
}

func (s *SpotifyService) BeginAuth() string {
	return s.spotify.BeginAuth()
}

func (s *SpotifyService) Callback(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.spotify.Callback(ctx, code)

}

func (s *SpotifyService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	return s.spotify.Refresh(ctx, refreshToken)
}
