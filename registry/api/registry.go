package api

import (
	"context"
	"net/http"
	"time"
	apiusecase "youtubelist/application/usecase/api"
	webusecase "youtubelist/application/usecase/web"
	"youtubelist/domain/config"
	"youtubelist/util/log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	"cloud.google.com/go/firestore"
	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/middleware"
)

func Start() {
	m := mux.NewRouter()
	ctx := context.Background()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	time.Local = loc
	fsCli, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		panic(err)
	}
	rd := redis.NewClient(&redis.Options{
		Addr:     "redis-14768.c1.asia-northeast1-1.gce.cloud.redislabs.com:14768",
		Password: "AcaU7b5eRS6x5YL9AOEYVAGkVX5mWrDd",
		DB:       0,
		PoolSize: 10,
	})

	// logger
	var logger log.Logger
	if config.IsLocal() {
		logger = log.NewlocalLogger()
	} else {
		logger = log.NewLogger()
		logz.SetConfig(logz.Config{
			ProjectID:      config.ProjectID,
			NeedsAccessLog: false,
		})
		logz.InitTracer()
	}
	apiusecase.RegisterUsecase(m, fsCli, logger, rd)
	webusecase.RegisterUsecase(m, fsCli, logger)
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
	logger.Criticalf(ctx, "%+v", server.ListenAndServe())
}
