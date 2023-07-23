package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	println("server is running:" + port)
	server := &http.Server{Handler: h, Addr: fmt.Sprintf(":%s", port)}
	base.Log.Criticalf(ctx, "%+v", server.ListenAndServe())
}
