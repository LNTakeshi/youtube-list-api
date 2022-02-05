package api

import (
	"context"
	"net/http"
	"time"
	apiusecase "youtubelist/application/usecase/api"
	webusecase "youtubelist/application/usecase/web"
	"youtubelist/domain/config"
	"youtubelist/util/log"

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

	// logger
	var logger log.Logger
	if config.IsLocal() {
		logger = log.NewlocalLogger()
	} else {
		logger = log.NewLogger()
		logz.SetConfig(logz.Config{
			ProjectID:      config.ProjectID,
			NeedsAccessLog: true, // Writes no access log
		})
		logz.InitTracer()
	}
	apiusecase.RegisterUsecase(m, fsCli, logger)
	webusecase.RegisterUsecase(m, fsCli, logger)
	m.HandleFunc("/hcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(m)
	println("server is running")
	server := &http.Server{Handler: h, Addr: ":8080"}
	logger.Criticalf(ctx, "%+v", server.ListenAndServe())
}
