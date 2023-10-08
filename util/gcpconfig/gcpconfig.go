package gcpconfig

import (
	"context"
	"os"
	"youtubelist/application/keystore"
	"youtubelist/domain/config"
)

type GcpConfig struct {
	RedisConfig   RedisConfig
	YoutubeConfig YoutubeConfig
	SpotifyConfig SpotifyConfig
}

type RedisConfig struct {
	Addr     string
	Password string
}

type YoutubeConfig struct {
	ApiKey string
}

type SpotifyConfig struct {
	ClientId     string
	ClientSecret string
}

func LoadGcpConfig(ctx context.Context, secretCli keystore.IKeyStore) *GcpConfig {
	println("islocal:", config.IsLocal())
	if config.IsLocal() {
		// TODO: ローカル環境の場合は未検討
		return &GcpConfig{
			RedisConfig: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
			},
		}
	}

	cfg := &GcpConfig{
		RedisConfig: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		YoutubeConfig: YoutubeConfig{
			ApiKey: os.Getenv("YOUTUBE_API_KEY"),
		},
		SpotifyConfig: SpotifyConfig{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		},
	}

	return cfg
}
