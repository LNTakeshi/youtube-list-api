package api

import (
	"encoding/json"
	"net/http"
	"youtubelist/application/service"
	"youtubelist/domain/config/constant"
	"youtubelist/errors"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"github.com/morikuni/failure"
)

type SendArgs struct {
	Username string `validate:"max=20"`
	Url      string `validate:"required,url"`
	Start    string
	End      string
	Title    string `validate:"max=100"`
	RoomID   string `schema:"room_id" validate:"required"`
	UUID     string `validate:"required" schema:"uuid"`
}

type Error struct {
	Error string `json:"error"`
}

var _ Writer = Error{}
var _ Writer = Success{}

type Success struct {
}

type Writer interface {
	WriteJSON(rw http.ResponseWriter)
}

func (e Error) WriteJSON(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
	b, _ := json.Marshal(e)
	rw.Write(b)
}

func (e Success) WriteJSON(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{}"))
}

func (u *Usecase) Send(rw http.ResponseWriter, r *http.Request) {
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	if r.Method != http.MethodPost {
		Success{}.WriteJSON(rw)
		// Error{Error: "リクエストがおかしい"}.WriteJSON(rw)
		// u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
	if err := r.ParseForm(); err != nil {
		Error{Error: "パースができない"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	args := &SendArgs{}
	if err := schema.NewDecoder().Decode(args, r.Form); err != nil {
		Error{Error: "パースができない"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	if err := validator.New().Struct(args); err != nil {
		if args.Url == "" {
			Error{Error: "URLが空"}.WriteJSON(rw)
			u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
			return
		}
		Error{Error: "バリデーションエラー。フォームに変なのが入っていそう"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	u.Log.Infof(ctx, "%+v", args)

	url, urlType, err := constant.NewUrlType(args.Url)
	if err != nil {
		Error{Error: "無効なyoutube/ニコニコ/SoundCloud/twitterのURLです"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}

	urlService := service.NewFetchUrlService(u.UsecaseBase)

	fetchResult, err := urlService.Fetch(ctx, urlType, url)
	if err != nil {
		Error{Error: "URLが取得できなかった"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	if args.Title != "" {
		if args.Title == "[HIDDEN]" {
			fetchResult.Title = args.Title + fetchResult.Title
		} else {
			fetchResult.Title = args.Title
		}
	}

	s := service.NewService(args.RoomID, args.UUID, u.FsCli, u.Redis)
	err = s.Add(ctx, fetchResult, args.Username, args.Start, args.End)
	if err != nil {
		if failure.Is(err, errors.ErrInvalidTime) {
			Error{Error: "Start/Endの指定が間違っている"}.WriteJSON(rw)
		} else {
			Error{Error: "URLが追加できなかった(内部エラー)"}.WriteJSON(rw)
		}
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	Success{}.WriteJSON(rw)
}
