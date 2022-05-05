package web

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"
	"youtubelist/errors"
	"youtubelist/react"
	"youtubelist/util/log"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/morikuni/failure"
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

type spaHandler struct {
	staticPath string
	indexPath  string
	staticFS   fs.FS
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
	_, err := h.staticFS.Open(r.URL.Path)
	if err != nil {
		// file does not exist, serve index.html
		index, err := h.staticFS.Open(h.indexPath)
		if err != nil {
			fmt.Printf("%+v\n", failure.New(errors.ErrBadUrl))

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		r, err := ioutil.ReadAll(index)
		if err != nil {
			fmt.Printf("%+v\n", failure.New(errors.ErrBadUrl))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(r)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.FS(h.staticFS)).ServeHTTP(w, r)
}

func RegisterUsecase(m *mux.Router, fsCli *firestore.Client, logger log.Logger) {
	// base := &UsecaseBase{FsCli: fsCli, Log: logger}

	// usecaseList := make([]*Usecase, 0)
	// usecase := &Usecase{FuncName: "/youtube-list/room/{RoomID}", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Room
	// usecaseList = append(usecaseList, usecase)

	// usecase = &Usecase{FuncName: "/youtube-list/{RoomID}", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Room
	// usecaseList = append(usecaseList, usecase)

	// usecase = &Usecase{FuncName: "/", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Index
	// usecaseList = append(usecaseList, usecase)

	// usecase = &Usecase{FuncName: "/youtube-list/", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Index
	// usecaseList = append(usecaseList, usecase)
	// for _, u := range usecaseList {
	// 	m.Path(u.FuncName).HandlerFunc(u.HandlerFunc).Methods("GET")
	// }

	spa := spaHandler{staticPath: "build", indexPath: "index.html", staticFS: react.Serve()}
	m.PathPrefix("/").Handler(spa)
}
