package util

import (
	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
	"youtubelist/application/niconico"
	"youtubelist/application/spotify"
	"youtubelist/application/youtube"
	yt_dlp "youtubelist/application/yt-dlp"
	"youtubelist/util/log"
)

type UsecaseBase struct {
	FsCli    *firestore.Client
	Log      log.Logger
	Redis    *redis.Client
	Youtube  youtube.IYoutube
	Niconico niconico.INiconico
	Spotify  spotify.ISpotify
	YtDlp    yt_dlp.IYtDlp
}
