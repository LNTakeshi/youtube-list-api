package twitter

import (
	"youtubelist/domain/entity"
)

type ITwitter interface {
	GetVideoInfo(url string) (*entity.FetchResult, error)
}
