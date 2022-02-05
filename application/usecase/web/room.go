package web

import (
	"embed"
	"net/http"
	"text/template"
	"youtubelist/errors"

	"github.com/gorilla/mux"
	"github.com/morikuni/failure"
)

//go:embed html
var f embed.FS

func (u *Usecase) Room(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["RoomID"]

	param := struct {
		RoomID string
	}{
		RoomID: id,
	}
	tmpl, err := template.ParseFS(f, "html/room.html")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}

	err = tmpl.Execute(rw, &param)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(""))
		u.Log.Errorf(ctx, "%+v", failure.New(errors.ErrBadRequest))
		return
	}
}
