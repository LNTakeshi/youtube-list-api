package service

import (
	"context"
	"youtubelist/application/usecase/util"
	"youtubelist/domain/config/constant"
	"youtubelist/domain/entity"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

type FetchUrlService struct {
	base *util.UsecaseBase
}

func NewFetchUrlService(base *util.UsecaseBase) *FetchUrlService {
	return &FetchUrlService{

		base: base,
	}
}

func (s *FetchUrlService) Fetch(ctx context.Context, urlType constant.UrlType, urlStr string) (*entity.FetchResult, error) {
	switch urlType {
	case constant.UrlTypeYoutube:
		result, err := s.base.Youtube.GetVideoInfo(urlStr)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return result, nil

	case constant.UrlTypeNicoNico:
		result, err := s.base.Niconico.GetVideoInfo(urlStr)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return result, nil
	case constant.UrlTypeTwitter:
		result, err := s.base.Twitter.GetVideoInfo(urlStr)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return result, nil

	case constant.UrlTypeSpotify:
		result, err := s.base.Spotify.GetAudioInfo(ctx, urlStr)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return result, nil
	}
	return nil, failure.New(errors.ErrFetchUrl)
}
