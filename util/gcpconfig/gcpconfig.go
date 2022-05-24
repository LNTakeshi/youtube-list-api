package gcpconfig

import (
	"context"
	"youtubelist/application/keystore"
	"youtubelist/domain/config"
	"youtubelist/domain/config/constant"
)

type GcpConfig struct {
	RedisConfig   RedisConfig
	YoutubeConfig YoutubeConfig
	TwitterConfig TwitterConfig
	SpotifyConfig SpotifyConfig
}

type RedisConfig struct {
	Addr     string
	Password string
}

type YoutubeConfig struct {
	ApiKey string
}

type TwitterConfig struct {
	ClientKey                 string
	ClientSecret              string
	ClientIdAccessToken       string
	ClientIdAccessTokenSecret string
}

type SpotifyConfig struct {
	ClientId     string
	ClientSecret string
}

func LoadGcpConfig(ctx context.Context, secretCli keystore.IKeyStore) *GcpConfig {
	if config.IsLocal() {
		// TODO: ローカル環境の場合は環境変数から取ってくるようにしたい
		return &GcpConfig{
			RedisConfig: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
			},
		}
	}
	redisAddr, err := secretCli.GetSecret(ctx, constant.SecretKeyNameRedisAddr)
	if err != nil {
		panic(err)
	}

	redisPassword, err := secretCli.GetSecret(ctx, constant.SecretKeyNameRedisPassword)
	if err != nil {
		panic(err)
	}

	youtubeApiKey, err := secretCli.GetSecret(ctx, constant.SecretKeyNameYoutubeApiKey)
	if err != nil {
		panic(err)
	}

	twitterClientKey, err := secretCli.GetSecret(ctx, constant.SecretKeyNameTwitterClientKey)
	if err != nil {
		panic(err)
	}

	twitterClientSecret, err := secretCli.GetSecret(ctx, constant.SecretKeyNameTwitterClientSecret)
	if err != nil {
		panic(err)
	}

	twitterClientIdAccessToken, err := secretCli.GetSecret(ctx, constant.SecretKeyNameTwitterClientIdAccessToken)
	if err != nil {
		panic(err)
	}

	twitterClientIdAccessTokenSecret, err := secretCli.GetSecret(ctx, constant.SecretKeyNameTwitterClientIdAccessTokenSecret)
	if err != nil {
		panic(err)
	}

	spotifyClientId, err := secretCli.GetSecret(ctx, constant.SecretKeyNameSpotifyClientId)
	if err != nil {
		panic(err)
	}

	spotifyClientSecret, err := secretCli.GetSecret(ctx, constant.SecretKeyNameSpotifyClientSecret)
	if err != nil {
		panic(err)
	}

	cfg := &GcpConfig{
		RedisConfig: RedisConfig{
			Addr:     string(redisAddr),
			Password: string(redisPassword),
		},
		YoutubeConfig: YoutubeConfig{
			ApiKey: string(youtubeApiKey),
		},
		TwitterConfig: TwitterConfig{
			ClientKey:                 string(twitterClientKey),
			ClientSecret:              string(twitterClientSecret),
			ClientIdAccessToken:       string(twitterClientIdAccessToken),
			ClientIdAccessTokenSecret: string(twitterClientIdAccessTokenSecret),
		},
		SpotifyConfig: SpotifyConfig{
			ClientId:     string(spotifyClientId),
			ClientSecret: string(spotifyClientSecret),
		},
	}

	return cfg
}
