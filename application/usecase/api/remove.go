package api

import (
	"net/http"
	"youtubelist/application/service"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"github.com/morikuni/failure"
)

type RemoveArgs struct {
	RoomID string `schema:"room_id" validate:"required"`
	UUID   string `schema:"uuid" validate:"required"`
	Index  string `schema:"index" validate:"required"`
}

func (u *Usecase) Remove(rw http.ResponseWriter, r *http.Request) {
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	ctx := r.Context()
	if r.Method != http.MethodPost {
		Success{}.WriteJSON(rw)
		return
	}
	if err := r.ParseForm(); err != nil {
		Error{Error: "パースができない"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	args := &RemoveArgs{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
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
	apiService := service.NewService(args.RoomID, args.UUID, u.FsCli, u.Redis)

	if err := apiService.Remove(ctx, args.Index); err != nil {
		Error{Error: "削除できなかった"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}

	Success{}.WriteJSON(rw)
}
