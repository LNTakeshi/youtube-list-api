package twitter

import (
	"context"
	"fmt"
	"github.com/morikuni/failure"
	"regexp"
	"strconv"
	"strings"
	"youtubelist/domain/entity"
	"youtubelist/errors"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitter2 "youtubelist/application/twitter"
	"youtubelist/util/gcpconfig"
)

type Client struct {
	cli *twitter.Client
}

func NewClient(ctx context.Context, config gcpconfig.TwitterConfig) twitter2.ITwitter {
	cfg := oauth1.NewConfig(config.ClientKey, config.ClientSecret)
	token := oauth1.NewToken(config.ClientIdAccessToken, config.ClientIdAccessTokenSecret)
	httpClient := cfg.Client(ctx, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	return &Client{
		cli: client,
	}
}

func (c Client) GetVideoInfo(url string) (*entity.FetchResult, error) {
	param := strings.Split(url, "/")
	id, err := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(param[len(param)-1]))
	if err != nil {
		return nil, failure.Wrap(err)
	}

	t, _, err := c.cli.Statuses.Show(int64(id), &twitter.StatusShowParams{})
	if err != nil {
		return nil, failure.Wrap(err)
	}
	var media *twitter.MediaEntity
	if t.ExtendedEntities != nil {
		for _, m := range t.ExtendedEntities.Media {
			if m.Type == "video" {
				media = &m
				break
			}
		}
	}
	if media == nil && t.Entities != nil {
		for _, m := range t.Entities.Media {
			if m.Type == "video" {
				media = &m
				break
			}
		}
	}
	if media == nil && t.ExtendedTweet != nil && t.ExtendedTweet.ExtendedEntities != nil {
		for _, m := range t.ExtendedTweet.ExtendedEntities.Media {
			if m.Type == "video" {
				media = &m
				break
			}
		}
	}
	if media == nil && t.ExtendedTweet != nil && t.ExtendedTweet.Entities != nil {
		for _, m := range t.ExtendedTweet.Entities.Media {
			if m.Type == "video" {
				media = &m
				break
			}
		}
	}
	if media == nil {
		return nil, failure.New(errors.ErrFetchUrl)
	}
	jsonUrl := fmt.Sprintf("https://twitter.com/%s/status/%d", t.User.ScreenName, id)
	length := media.VideoInfo.DurationMillis / 1000
	title := t.Text
	return &entity.FetchResult{
		Title:  title,
		Url:    jsonUrl,
		Length: length,
	}, nil
}
