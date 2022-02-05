package api

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
	usecase := &Usecase{FuncName: "/api/youtubelist/getList", UsecaseBase: base}
	usecase.HandlerFunc = usecase.GetList
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/api/youtubelist/send", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Send
	usecaseList = append(usecaseList, usecase)
	usecase = &Usecase{FuncName: "/api/youtubelist/remove", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Remove
	usecaseList = append(usecaseList, usecase)
	for _, u := range usecaseList {
		m.Handle(u.FuncName, u.HandlerFunc)
	}
}
