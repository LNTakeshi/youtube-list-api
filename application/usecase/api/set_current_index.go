package api

import (
	"net/http"
	"youtubelist/application/service"
	"youtubelist/errors"

	"github.com/glassonion1/logz"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"github.com/morikuni/failure"
)

type SetCurrentIndexArgs struct {
	RoomID   string `schema:"room_id" validate:"required"`
	MasterID string `schema:"master_id" validate:"required"`
	Index    string `schema:"index" validate:"required"`
}

func (u *Usecase) SetCurrentIndex(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		Error{Error: "リクエストがおかしい"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
	if err := r.ParseForm(); err != nil {
		Error{Error: "パースができない"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	args := &SetCurrentIndexArgs{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	logz.Infof(ctx, "SetCurrentIndexFromUnityClient: %+v", r.Form)
	if err := decoder.Decode(args, r.Form); err != nil {
		Error{Error: "パースができない"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	if err := validator.New().Struct(args); err != nil {
		Error{Error: "バリデーションエラー"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	apiService := service.NewService(args.RoomID, args.MasterID, u.FsCli, u.Redis)

	if err := apiService.SetCurrentIndex(ctx, args.Index); err != nil {
		Error{Error: "削除できなかった"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}

	Success{}.WriteJSON(rw)
}
