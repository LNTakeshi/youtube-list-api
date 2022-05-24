package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"youtubelist/domain/config/constant"
	"youtubelist/errors"

	"youtubelist/util/log"

	"github.com/morikuni/failure"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyService struct {
	log log.Logger
}

func NewSpotifyService(log log.Logger) *SpotifyService {
	return &SpotifyService{log: log}
}

func (s *SpotifyService) BeginAuth(ctx context.Context) string {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(constant.SPOTIFY_CLIENT_ID),
		spotifyauth.WithClientSecret(constant.SPOTIFY_CLIENT_SECRET),
		spotifyauth.WithRedirectURL("https://lntk.info/spotifycallback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserModifyPlaybackState))
	return auth.AuthURL("hoge")
}

func (s *SpotifyService) Callback(ctx context.Context, code string) (*oauth2.Token, error) {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(constant.SPOTIFY_CLIENT_ID),
		spotifyauth.WithClientSecret(constant.SPOTIFY_CLIENT_SECRET),
		spotifyauth.WithRedirectURL("https://lntk.info/spotifycallback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserModifyPlaybackState))
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return token, nil
}

func (s *SpotifyService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	cli := http.DefaultClient
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {strings.TrimSpace(refreshToken)},
	}
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", failure.Wrap(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(constant.SPOTIFY_CLIENT_ID, constant.SPOTIFY_CLIENT_SECRET)
	res, err := cli.Do(req)
	if err != nil {
		return "", failure.Wrap(err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		s.log.Errorf(ctx, "spotify refresh token failed. body: %s", b)
		return "", failure.New(errors.ErrBadRequest)
	}
	result := struct {
		AccessToken string `json:"access_token"`
	}{}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", failure.Wrap(err)
	}
	if err = json.Unmarshal(b, &result); err != nil {
		return "", failure.Wrap(err)
	}

	return result.AccessToken, nil
}
