package api

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/glassonion1/logz"
	"github.com/go-redis/redis/v8"
	"youtubelist/application/niconico"
	"youtubelist/application/spotify"
	"youtubelist/application/youtube"
	yt_dlp "youtubelist/application/yt-dlp"
	"youtubelist/domain/config"
	infraNiconico "youtubelist/infra/niconico"
	infraSpotify "youtubelist/infra/spotify"
	infraYoutube "youtubelist/infra/youtube"
	infraYtDlp "youtubelist/infra/yt-dlp"
	"youtubelist/util/gcpconfig"
	"youtubelist/util/log"
)

func provideFireStoreClient(ctx context.Context) *firestore.Client {
	fsCli, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		panic(err)
	}
	return fsCli
}

func provideLogger() log.Logger {
	var logger log.Logger
	if config.IsLocal() {
		logger = log.NewlocalLogger()
	} else {
		logger = log.NewLogger()
		logz.SetConfig(logz.Config{
			ProjectID:      config.ProjectID,
			NeedsAccessLog: false,
		})
		logz.InitTracer()
	}
	return logger
}

func provideRedisClient(cfg *gcpconfig.GcpConfig) *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Addr,
		Password: cfg.RedisConfig.Password,
	})
	return rd
}

func provideNiconico() niconico.INiconico {
	return infraNiconico.NewClient()
}

func provideYoutube(ctx context.Context, cfg *gcpconfig.GcpConfig) youtube.IYoutube {
	return infraYoutube.NewClient(ctx, cfg.YoutubeConfig)
}

func provideSpotify(cfg *gcpconfig.GcpConfig, logger log.Logger) spotify.ISpotify {
	return infraSpotify.NewClient(cfg.SpotifyConfig, logger)
}

func provideYtDlp() yt_dlp.IYtDlp {
	return infraYtDlp.NewClient()
}
