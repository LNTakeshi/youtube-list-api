package constant

import (
	"strings"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

const YoutubeApiKey = "AIzaSyBgt6J1j38WWdBHRzWQmF65qbZk2ltEygw"
const TWITTER_CLIENT_KEY = "grHV1oz4Y3zhznSmfxEdg"
const TWITTER_CLIENT_SECRET = "BD6Fj3A2MHdcd13iblzUfzrYYsY3leysqhS2gDrSs"
const TWITTER_CLIENT_ID_ACCESS_TOKEN = "63342875-Yc114gVRVyWsdz06FeGaI2frsjhN5gTdb3WvSpGAV"
const TWITTER_CLIENT_ID_ACCESS_TOKEN_SECRET = "MSY3lnc095Vu2qH3R9WIOej6ayJbMEDdFnMT7p9VviE4I"
const SPOTIFY_CLIENT_ID = "f4801f15efcd49e3bebf00e086cddfed"
const SPOTIFY_CLIENT_SECRET = "2f4fb3692c5d432d9c156f1e9ce53b63"

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
	case strings.HasPrefix(Url, "https://youtube.com/"), strings.HasPrefix(Url, "https://www.youtube.com/"), strings.HasPrefix(Url, "https://m.youtube.com/"):
		return Url, UrlTypeYoutube, nil
	case strings.HasPrefix(Url, "https://nico.ms/sm"), strings.HasPrefix(Url, "https://www.nicovideo.jp/watch/sm"), strings.HasPrefix(Url, "https://sp.nicovideo.jp/watch/sm"):
		return Url, UrlTypeNicoNico, nil
	case strings.HasPrefix(Url, "https://soundcloud.com/"):
		return Url, UrlTypeSoundCloud, nil
	case strings.HasPrefix(Url, "https://twitter.com/"):
		return Url, UrlTypeTwitter, nil
	case strings.HasPrefix(Url, "https://open.spotify.com/track/"):
		return Url, UrlTypeSpotify, nil
		// return "", UrlTypeUnknown, failure.New(errors.ErrBadUrl)
	}
	return "", UrlTypeUnknown, failure.New(errors.ErrBadUrl)
}
