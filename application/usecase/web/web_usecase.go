package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/morikuni/failure"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"youtubelist/application/usecase/util"
	"youtubelist/errors"
	"youtubelist/react"
)

type Usecase struct {
	FuncName    string
	HandlerFunc http.HandlerFunc
	*util.UsecaseBase
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
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/ytl/")
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
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.FileServer(http.FS(h.staticFS)).ServeHTTP(w, r)
}

func RegisterUsecase(m *mux.Router, base *util.UsecaseBase) {
	usecaseList := make([]*Usecase, 0)
	// usecase := &Usecase{FuncName: "/youtube-list/room/{RoomID}", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Room
	// usecaseList = append(usecaseList, usecase)

	// usecase = &Usecase{FuncName: "/youtube-list/{RoomID}", UsecaseBase: base}
	// usecase.HandlerFunc = usecase.Room
	// usecaseList = append(usecaseList, usecase)

	usecase := &Usecase{FuncName: "/spotifycallback", UsecaseBase: base}
	usecase.HandlerFunc = usecase.SpotifyCallback
	usecaseList = append(usecaseList, usecase)

	usecase = &Usecase{FuncName: "/{Num}", UsecaseBase: base}
	usecase.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["Num"]

		_, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Location", "/ytl/"+id)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
	usecaseList = append(usecaseList, usecase)

	usecase = &Usecase{FuncName: "/youtube-list/myip", UsecaseBase: base}
	usecase.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		s := strings.Split(r.RemoteAddr, ":")
		w.Write([]byte(fmt.Sprintf("<html><body><div class=\"myip\" id=\"%s\"></div></body></html>", s[0])))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
	usecaseList = append(usecaseList, usecase)

	usecase = &Usecase{FuncName: "/", UsecaseBase: base}
	usecase.HandlerFunc = usecase.Index
	usecaseList = append(usecaseList, usecase)
	for _, u := range usecaseList {
		m.Path(u.FuncName).HandlerFunc(u.HandlerFunc).Methods("GET")
	}

	spa := spaHandler{staticPath: "build", indexPath: "index.html", staticFS: react.Serve()}
	m.PathPrefix("/ytl").Handler(spa)
}
