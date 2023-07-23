package constant

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

type UrlType int

const (
	UrlTypeUnknown UrlType = iota
	UrlTypeYoutube
	UrlTypeNicoNico
	UrlTypeSoundCloud
	UrlTypeTwitter
	UrlTypeSpotify
)

func NewUrlType(Url string) (string, UrlType, error) {
	if strings.HasPrefix(Url, "https://youtu.be/") {
		Url = "https://www.youtube.com/watch?v=" + strings.TrimPrefix(Url, "https://youtu.be/")
	}
	switch {
	case strings.HasPrefix(Url, "https://www.youtube.com/shorts/"):
		u, err := url.Parse(Url)
		if err != nil {
			return "", UrlTypeUnknown, failure.New(errors.ErrBadUrl)
		}
		_, p := path.Split(u.Path)
		return fmt.Sprintf("https://www.youtube.com/watch?v=%s", p), UrlTypeYoutube, nil
	case strings.HasPrefix(Url, "https://youtube.com/"), strings.HasPrefix(Url, "https://www.youtube.com/"), strings.HasPrefix(Url, "https://m.youtube.com/"):
		return strings.Split(Url, "&")[0], UrlTypeYoutube, nil
	case strings.HasPrefix(Url, "https://nico.ms/sm"), strings.HasPrefix(Url, "https://www.nicovideo.jp/watch/sm"), strings.HasPrefix(Url, "https://sp.nicovideo.jp/watch/sm"):
		return Url, UrlTypeNicoNico, nil
	case strings.HasPrefix(Url, "https://soundcloud.com/"):
		return Url, UrlTypeSoundCloud, nil
	case strings.HasPrefix(Url, "https://twitter.com/"):
		return Url, UrlTypeTwitter, nil
	case strings.HasPrefix(Url, "https://open.spotify.com/"):
		return Url, UrlTypeSpotify, nil
	}
	return "", UrlTypeUnknown, failure.New(errors.ErrBadUrl)
}

type SecretKeyName string

const (
	SecretKeyNameYoutubeApiKey       SecretKeyName = "youtube-api-key"
	SecretKeyNameSpotifyClientId                   = "spotify-client-id"
	SecretKeyNameSpotifyClientSecret               = "spotify-client-secret"
	SecretKeyNameRedisAddr                         = "redis-addr"
	SecretKeyNameRedisPassword                     = "redis-password"
)
