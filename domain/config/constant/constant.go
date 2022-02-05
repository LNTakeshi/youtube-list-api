package constant

import (
	"strings"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

const YoutubeApiKey = "AIzaSyBgt6J1j38WWdBHRzWQmF65qbZk2ltEygw"

type UrlType int

const (
	UrlTypeUnknown UrlType = iota
	UrlTypeYoutube
	UrlTypeNicoNico
	UrlTypeSoundCloud
	UrlTypeTwitter
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
	}
	return "", UrlTypeUnknown, failure.New(errors.ErrBadUrl)
}
