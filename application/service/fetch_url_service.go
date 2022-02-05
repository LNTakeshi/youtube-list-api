package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"youtubelist/domain/config/constant"
	"youtubelist/errors"

	"github.com/morikuni/failure"
	"github.com/rickb777/date/period"
)

type FetchUrlService struct {
	UrlType constant.UrlType
	Url     string
}

type FetchResult struct {
	Url    string
	Title  string
	Length int
}

func NewFetchUrlService(urlType constant.UrlType, url string) *FetchUrlService {
	return &FetchUrlService{
		UrlType: urlType,
		Url:     url,
	}
}

func (s *FetchUrlService) Fetch(ctx context.Context) (*FetchResult, error) {
	cli := http.DefaultClient
	switch s.UrlType {
	case constant.UrlTypeYoutube:
		apiUrl, err := url.Parse(s.Url)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		if len(apiUrl.Query().Get("v")) == 0 {
			return nil, failure.New(errors.ErrFetchUrl)
		}
		u, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
		if err != nil {
			return nil, failure.Wrap(err)
		}
		q := url.Values{}
		q.Set("id", apiUrl.Query().Get("v"))
		q.Set("key", constant.YoutubeApiKey)
		q.Set("part", "snippet,contentDetails")
		u.RawQuery = q.Encode()
		println(u.String())
		res, err := cli.Get(u.String())
		if err != nil {
			return nil, failure.Wrap(err)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		js := struct {
			Items []struct {
				Snippet struct {
					Title string `json:"title"`
				} `json:"snippet"`
				ContentDetails struct {
					Duration string `json:"duration"`
				} `json:"contentDetails"`
			} `json:"items"`
		}{}
		err = json.Unmarshal(b, &js)
		if err != nil || len(js.Items) == 0 {
			return nil, failure.Wrap(err)
		}
		fmt.Printf("%+v\n", js)
		pe, err := period.Parse(js.Items[0].ContentDetails.Duration)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		du, ok := pe.Duration()
		if !ok {
			return nil, failure.New(errors.ErrFetchUrl)
		}
		return &FetchResult{
			Title:  js.Items[0].Snippet.Title,
			Url:    s.Url,
			Length: int(du.Seconds()),
		}, nil
	}
	return nil, failure.New(errors.ErrFetchUrl)
}
