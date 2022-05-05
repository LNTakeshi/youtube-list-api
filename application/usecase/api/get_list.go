package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"youtubelist/application/service"
	"youtubelist/application/status"
	"youtubelist/errors"

	"github.com/google/uuid"
	"github.com/morikuni/failure"
)

type GetListArgs struct {
	RoomID         string
	UUID           string
	MasterID       string
	LastUpdateDate time.Time
}

func (u *Usecase) GetList(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("リクエストがおかしい"))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
	if err := r.ParseForm(); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("パースができない"))
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	args := &GetListArgs{}
	args.RoomID = r.Form.Get("room_id")
	if args.RoomID == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("room_idがありません"))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
	args.UUID = r.Form.Get("uuid")
	if args.UUID == "" {
		args.UUID = uuid.NewString()
	}
	args.MasterID = r.Form.Get("master_id")
	args.LastUpdateDate, _ = time.Parse(time.RFC3339, r.Form.Get("lastUpdateDate"))
	fmt.Printf("aa:%+v", args.LastUpdateDate)

	apiService := service.NewService(args.RoomID, args.UUID, u.FsCli, u.Redis)

	getList, uuid, err := apiService.GetRoom(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("部屋が作れなかった"))
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}

	st := status.NewGetList(getList, uuid, args.LastUpdateDate)
	res, err := json.Marshal(st)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(res)
}
