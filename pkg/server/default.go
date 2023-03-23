package server

import (
	"log"
	"net/http"
	"time"
)

func HTTPStandardServer(addr string, h http.Handler) error {
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("Starting default server on %s", addr)
	return srv.ListenAndServe()
}
