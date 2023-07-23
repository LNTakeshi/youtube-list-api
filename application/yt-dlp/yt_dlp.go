package yt_dlp

import (
	"context"
	"youtubelist/domain/entity"
)

type IYtDlp interface {
	GetVideoInfo(ctx context.Context, requestURL string) (*entity.FetchResult, error)
}
