package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"youtubelist/application/service"

	"github.com/morikuni/failure"
)

func (u *Usecase) SpotifyRefresh(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	requestUrl, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		Error{Error: "URLパースエラー"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}

	u.Log.Debugf(ctx, "SpotifyRefresh: %+v", requestUrl.Query())
	s := service.NewSpotifyService(u.UsecaseBase)
	token, err := s.Refresh(ctx, requestUrl.Query().Get("token"))
	if err != nil {
		Error{Error: "リフレッシュ失敗"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	rw.WriteHeader(http.StatusOK)
	b, err := json.Marshal(struct{ Code string }{Code: token})
	if err != nil {
		Error{Error: "parse error"}.WriteJSON(rw)
		u.Log.Errorf(ctx, "%+v", failure.Wrap(err))
		return
	}
	rw.Write(b)
}
