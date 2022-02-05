package util

import (
	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
	"youtubelist/application/niconico"
	"youtubelist/application/spotify"
	"youtubelist/application/twitter"
	"youtubelist/application/youtube"
	"youtubelist/util/log"
)

type UsecaseBase struct {
	FsCli    *firestore.Client
	Log      log.Logger
	Redis    *redis.Client
	Twitter  twitter.ITwitter
	Youtube  youtube.IYoutube
	Niconico niconico.INiconico
	Spotify  spotify.ISpotify
}
