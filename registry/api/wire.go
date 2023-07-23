//go:build wireinject
// +build wireinject

package api

import (
	"context"
	"github.com/google/wire"
	"youtubelist/application/usecase/api"
	"youtubelist/application/usecase/util"
	"youtubelist/util/gcpconfig"
)

func InitUsecaseBase(ctx context.Context, cfg *gcpconfig.GcpConfig) *util.UsecaseBase {
	wire.Build(api.NewUsecaseBase,
		provideFireStoreClient,
		provideLogger,
		provideRedisClient,
		provideNiconico,
		provideYoutube,
		provideSpotify,
		provideYtDlp)
	return nil
}
