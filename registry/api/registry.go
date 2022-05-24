package api

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	apiusecase "youtubelist/application/usecase/api"
	webusecase "youtubelist/application/usecase/web"
	"youtubelist/infra/keystore/secretmanager"
	"youtubelist/util/gcpconfig"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/middleware"
)

func Start() {
	m := mux.NewRouter()
	ctx := context.Background()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	time.Local = loc

	secret, err := secretmanager.NewClient(ctx)
	if err != nil {
		logz.Errorf(ctx, "%+v", err)
		return
	}

	gcpConfig := gcpconfig.LoadGcpConfig(ctx, secret)

	base := InitUsecaseBase(ctx, gcpConfig)

	apiusecase.RegisterUsecase(m, base)
	webusecase.RegisterUsecase(m, base)
	m.HandleFunc("/{RoomID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["RoomID"]
		http.Redirect(w, r, "/youtube-list/room/"+id, http.StatusMovedPermanently)
	}).Methods("GET")
	m.HandleFunc("/hcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(m)

	println("server is running")
	server := &http.Server{Handler: h, Addr: ":8080"}
	base.Log.Criticalf(ctx, "%+v", server.ListenAndServe())
}
