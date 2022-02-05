package spotify

import (
	"context"
	"golang.org/x/oauth2"
	"youtubelist/domain/entity"
)

type ISpotify interface {
	GetAudioInfo(ctx context.Context, url string) (*entity.FetchResult, error)
	BeginAuth() string
	Callback(ctx context.Context, code string) (*oauth2.Token, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
}
