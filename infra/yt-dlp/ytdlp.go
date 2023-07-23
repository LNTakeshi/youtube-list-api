package yt_dlp

import (
	"context"
	"github.com/morikuni/failure"
	"github.com/wader/goutubedl"
	"youtubelist/domain/entity"
)

type Client struct {
}

func NewClient() *Client {
	goutubedl.Path = "yt-dlp"

	return &Client{}
}

func (c *Client) GetVideoInfo(ctx context.Context, requestURL string) (*entity.FetchResult, error) {
	result, err := goutubedl.New(ctx, requestURL, goutubedl.Options{})
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return &entity.FetchResult{
		Url:    requestURL,
		Title:  result.Info.Title,
		Length: int(result.Info.Duration),
	}, nil
}
