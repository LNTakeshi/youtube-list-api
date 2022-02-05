package niconico

import (
	"youtubelist/domain/entity"
)

type INiconico interface {
	GetVideoInfo(url string) (*entity.FetchResult, error)
}
