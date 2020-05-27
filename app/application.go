package app

import (
	"github.com/gorilla/mux"
	"github.com/willqiang/bookstore_items-api/clients/elasticsearch"
	"net/http"
	"time"
)

var (
	router = mux.NewRouter()
)

func StartApplication() {
	elasticsearch.Init()

	mapUrls()

	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      500 * time.Millisecond,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           router,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
