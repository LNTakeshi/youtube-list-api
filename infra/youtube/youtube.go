package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/morikuni/failure"
	"github.com/rickb777/date/period"
	"io"
	"net/http"
	"net/url"
	"youtubelist/application/youtube"
	"youtubelist/domain/entity"
	"youtubelist/errors"

	"youtubelist/util/gcpconfig"
)

type Client struct {
	cli    *http.Client
	config gcpconfig.YoutubeConfig
}

func NewClient(ctx context.Context, config gcpconfig.YoutubeConfig) youtube.IYoutube {
	return &Client{
		cli:    http.DefaultClient,
		config: config,
	}
}

func (c Client) GetVideoInfo(urlString string) (*entity.FetchResult, error) {
	apiUrl, err := url.Parse(urlString)
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
	q.Set("key", c.config.ApiKey)
	q.Set("part", "snippet,contentDetails")
	u.RawQuery = q.Encode()
	res, err := c.cli.Get(u.String())
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
	urlString = "https://www.youtube.com/watch?v=" + apiUrl.Query().Get("v")
	return &entity.FetchResult{
		Title:  js.Items[0].Snippet.Title,
		Url:    urlString,
		Length: int(du.Seconds()),
	}, nil
}
