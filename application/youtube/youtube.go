package youtube

import (
	"youtubelist/domain/entity"
)

type IYoutube interface {
	GetVideoInfo(url string) (*entity.FetchResult, error)
}
