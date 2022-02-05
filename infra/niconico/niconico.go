package niconico

import (
	"encoding/xml"
	"github.com/morikuni/failure"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"youtubelist/application/niconico"
	"youtubelist/domain/entity"
	"youtubelist/errors"
)

type Client struct {
	cli *http.Client
}

func NewClient() niconico.INiconico {
	return &Client{
		cli: http.DefaultClient,
	}
}

func (c Client) GetVideoInfo(urlStr string) (*entity.FetchResult, error) {
	urlStr = strings.TrimPrefix(urlStr, "https://www.nicovideo.jp/watch/sm")
	urlStr = strings.TrimPrefix(urlStr, "https://sp.nicovideo.jp/watch/sm")
	urlStr = strings.TrimPrefix(urlStr, "https://nico.ms/sm")
	id := regexp.MustCompile(`[0-9]+`).FindString(urlStr)
	res, err := c.cli.Get("https://ext.nicovideo.jp/api/getthumbinfo/sm" + id)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	x := struct {
		Thumb struct {
			Title  string `xml:"title"`
			Length string `xml:"length"`
			UserID int    `xml:"user_id"`
		} `xml:"thumb"`
	}{}
	err = xml.Unmarshal(b, &x)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	if x.Thumb.UserID == 30338847 {
		return nil, failure.New(errors.ErrFetchUrl)
	}

	jsonUrl := "https://www.nicovideo.jp/watch/sm" + id
	if x.Thumb.Title == "" {
		return nil, failure.New(errors.ErrAdd)
	}
	t := strings.Split(x.Thumb.Length, ":")
	if len(t) != 2 {
		return nil, failure.New(errors.ErrFetchUrl)
	}
	m, err := strconv.Atoi(t[0])
	if err != nil {
		return nil, failure.Wrap(err)
	}
	s, err := strconv.Atoi(t[1])
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return &entity.FetchResult{
		Title:  x.Thumb.Title,
		Url:    jsonUrl,
		Length: m*60 + s,
	}, nil
}
