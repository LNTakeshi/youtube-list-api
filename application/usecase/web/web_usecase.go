package web

import (
	"net/http"
	"youtubelist/util/log"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

type Usecase struct {
	FuncName    string
	HandlerFunc http.HandlerFunc
	*UsecaseBase
}
type UsecaseBase struct {
	FsCli *firestore.Client
	Log   log.Logger
}

func RegisterUsecase(m *mux.Router, fsCli *firestore.Client, logger log.Logger) {
	base := &UsecaseBase{FsCli: fsCli, Log: logger}

	usecaseList := make([]*Usecase, 0)
	usecase := &Usecase{FuncName: "/youtube-list/room/{RoomID}", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Room
	usecaseList = append(usecaseList, usecase)

	usecase = &Usecase{FuncName: "/youtube-list/{RoomID}", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Room
	usecaseList = append(usecaseList, usecase)
	for _, u := range usecaseList {
		m.Path(u.FuncName).HandlerFunc(u.HandlerFunc).Methods("GET")
	}
}
