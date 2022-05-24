package service

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"youtubelist/domain/config/constant"
	"youtubelist/errors"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/morikuni/failure"
	"github.com/rickb777/date/period"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
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
		s.Url = "https://www.youtube.com/watch?v=" + apiUrl.Query().Get("v")
		return &FetchResult{
			Title:  js.Items[0].Snippet.Title,
			Url:    s.Url,
			Length: int(du.Seconds()),
		}, nil
	case constant.UrlTypeNicoNico:
		url := strings.TrimPrefix(s.Url, "https://www.nicovideo.jp/watch/sm")
		url = strings.TrimPrefix(url, "https://sp.nicovideo.jp/watch/sm")
		url = strings.TrimPrefix(url, "https://nico.ms/sm")
		id := regexp.MustCompile(`[0-9]+`).FindString(url)
		res, err := cli.Get("https://ext.nicovideo.jp/api/getthumbinfo/sm" + id)
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
		return &FetchResult{
			Title:  x.Thumb.Title,
			Url:    jsonUrl,
			Length: m*60 + s,
		}, nil
	case constant.UrlTypeTwitter:
		param := strings.Split(s.Url, "/")
		id, err := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(param[len(param)-1]))
		if err != nil {
			return nil, failure.Wrap(err)
		}

		config := oauth1.NewConfig(constant.TWITTER_CLIENT_KEY, constant.TWITTER_CLIENT_SECRET)
		token := oauth1.NewToken(constant.TWITTER_CLIENT_ID_ACCESS_TOKEN, constant.TWITTER_CLIENT_ID_ACCESS_TOKEN_SECRET)
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)
		t, _, err := client.Statuses.Show(int64(id), &twitter.StatusShowParams{})
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
		return &FetchResult{
			Title:  title,
			Url:    jsonUrl,
			Length: length,
		}, nil
	case constant.UrlTypeSpotify:
		apiUrl, err := url.Parse(s.Url)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		sp := strings.Split(apiUrl.Path, "/")
		id := spotify.ID(sp[len(sp)-1])
		ctx := context.Background()
		config := &clientcredentials.Config{
			ClientID:     constant.SPOTIFY_CLIENT_ID,
			ClientSecret: constant.SPOTIFY_CLIENT_SECRET,
			TokenURL:     spotifyauth.TokenURL,
		}
		token, err := config.Token(ctx)
		if err != nil {
			log.Fatalf("couldn't get token: %v", err)
		}
		httpClient := spotifyauth.New().Client(ctx, token)
		client := spotify.New(httpClient)
		track, err := client.GetTrack(ctx, id)
		if err != nil {
			log.Fatalf("couldn't get token: %v", err)
		}
		artistName := ""
		if len(track.Artists) > 0 {
			artistName = track.Artists[0].Name
		}
		return &FetchResult{
			Title:  fmt.Sprintf("%s / %s", track.Name, artistName),
			Url:    s.Url,
			Length: int(time.Duration(track.TimeDuration().Seconds())),
		}, nil
	}
	return nil, failure.New(errors.ErrFetchUrl)
}
