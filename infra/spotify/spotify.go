package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/morikuni/failure"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"youtubelist/util/log"

	"net/http"
	"net/url"
	"strings"
	"time"
	iSpotify "youtubelist/application/spotify"
	"youtubelist/domain/entity"
	"youtubelist/errors"
	"youtubelist/util/gcpconfig"
)

type Client struct {
	config *clientcredentials.Config
	log    log.Logger
}

func NewClient(config gcpconfig.SpotifyConfig, log log.Logger) iSpotify.ISpotify {
	return &Client{
		config: &clientcredentials.Config{
			ClientID:     config.ClientId,
			ClientSecret: config.ClientSecret,
			TokenURL:     spotifyauth.TokenURL,
		},
		log: log,
	}
}

func (c *Client) GenerateClient(ctx context.Context) (*spotify.Client, error) {
	token, err := c.config.Token(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client, nil
}

func (c *Client) GetAudioInfo(ctx context.Context, urlStr string) (*entity.FetchResult, error) {
	apiUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	sp := strings.Split(apiUrl.Path, "/")
	id := spotify.ID(sp[len(sp)-1])

	client, err := c.GenerateClient(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	track, err := client.GetTrack(ctx, id)
	if err != nil {
		c.log.Criticalf(ctx, "couldn't get token: %v", err)
	}
	artistName := ""
	if len(track.Artists) > 0 {
		artistName = track.Artists[0].Name
	}
	return &entity.FetchResult{
		Title:  fmt.Sprintf("%s / %s", track.Name, artistName),
		Url:    urlStr,
		Length: int(time.Duration(track.TimeDuration().Seconds())),
	}, nil
}

func (c *Client) BeginAuth() string {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(c.config.ClientID),
		spotifyauth.WithClientSecret(c.config.ClientSecret),
		spotifyauth.WithRedirectURL("https://lntk.info/spotifycallback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserModifyPlaybackState))
	return auth.AuthURL("hoge")
}

func (c *Client) Callback(ctx context.Context, code string) (*oauth2.Token, error) {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(c.config.ClientID),
		spotifyauth.WithClientSecret(c.config.ClientSecret),
		spotifyauth.WithRedirectURL("https://lntk.info/spotifycallback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserModifyPlaybackState))
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return token, nil
}

func (c *Client) Refresh(ctx context.Context, refreshToken string) (string, error) {
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
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)
	res, err := cli.Do(req)
	if err != nil {
		return "", failure.Wrap(err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		c.log.Errorf(ctx, "spotify refresh token failed. body: %s", b)
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
