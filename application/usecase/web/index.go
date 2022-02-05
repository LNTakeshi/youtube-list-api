package web

import (
	"net/http"
	"text/template"
	"youtubelist/errors"

	"github.com/morikuni/failure"
)

func (u *Usecase) Index(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tmpl, err := template.ParseFS(f, "html/index.html")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}

	err = tmpl.Execute(rw, struct{}{})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
}
