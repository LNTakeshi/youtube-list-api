package api

import (
	"net/http"
	"youtubelist/errors"

	"github.com/glassonion1/logz"
	"github.com/morikuni/failure"
)

func (u *Usecase) SendError(rw http.ResponseWriter, r *http.Request) {
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
	logz.Errorf(ctx, "SendErrorFromUnityClient: %+v", r.Form)

	Success{}.WriteJSON(rw)
}
