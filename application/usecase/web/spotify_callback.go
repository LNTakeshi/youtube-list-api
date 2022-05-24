package web

import (
	"net/http"
	"net/url"
	"text/template"
	"youtubelist/application/service"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

func (u *Usecase) SpotifyCallback(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// rw.Header().Set("Access-Control-Allow-Headers", "*")
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	// rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	s := service.NewSpotifyService(u.UsecaseBase)

	requestUrl, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
	token, err := s.Callback(ctx, requestUrl.Query()["code"][0])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}

	tmpl, err := template.ParseFS(f, "html/spotify_callback.html")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}

	err = tmpl.Execute(rw, &token)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
}
