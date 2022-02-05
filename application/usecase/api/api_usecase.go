package api

import (
	"net/http"
	"youtubelist/application/niconico"
	"youtubelist/application/spotify"
	"youtubelist/application/twitter"
	"youtubelist/application/usecase/util"
	"youtubelist/application/youtube"
	"youtubelist/util/log"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Usecase struct {
	FuncName    string
	HandlerFunc http.HandlerFunc
	*util.UsecaseBase
}

func NewUsecaseBase(
	fsCli *firestore.Client,
	logger log.Logger,
	rd *redis.Client,
	twitter twitter.ITwitter,
	youtube youtube.IYoutube,
	niconico niconico.INiconico,
	iSpotify spotify.ISpotify,
) *util.UsecaseBase {
	return &util.UsecaseBase{
		FsCli:    fsCli,
		Log:      logger,
		Redis:    rd,
		Twitter:  twitter,
		Youtube:  youtube,
		Niconico: niconico,
		Spotify:  iSpotify,
	}
}

func RegisterUsecase(m *mux.Router, base *util.UsecaseBase) {
	usecaseList := make([]*Usecase, 0)
	usecase := &Usecase{FuncName: "/youtube-list/api/youtubelist/getList", UsecaseBase: base}
	usecase.HandlerFunc = usecase.GetList
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/send", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Send
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/remove", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Remove
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/setCurrentIndex", UsecaseBase: base}
	usecase.HandlerFunc = usecase.SetCurrentIndex
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/sendError", UsecaseBase: base}
	usecase.HandlerFunc = usecase.SendError
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/spotifyAuth", UsecaseBase: base}
	usecase.HandlerFunc = usecase.SpotifyAuth
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/youtube-list/api/youtubelist/spotifyRefresh", UsecaseBase: base}
	usecase.HandlerFunc = usecase.SpotifyRefresh
	usecaseList = append(usecaseList, usecase)
	for _, u := range usecaseList {
		m.Handle(u.FuncName, u.HandlerFunc)
	}
}
