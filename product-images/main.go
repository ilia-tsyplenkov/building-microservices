package main

import (
	"building-microservices/product-images/files"
	"building-microservices/product-images/handlers"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

func main() {
	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString("DEBUG"),
		},
	)

	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})
	basePath := "./imagestore"
	store, err := files.NewLocal(basePath, 5*1024*1024)
	if err != nil {
		l.Error("unable to create store", "err", err)
		os.Exit(1)
	}
	fh := handlers.NewFiles(l, store)
	mwh := handlers.NewGzipHandler()

	r := mux.NewRouter()
	postFiles := r.Methods(http.MethodPost).Subrouter()
	postFiles.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{1,4}}", fh)

	getFiles := r.Methods(http.MethodGet).Subrouter()
	getFiles.Use(mwh.GzipMiddleWare)
	getFiles.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{1,4}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(basePath))))

	server := &http.Server{
		Addr:         ":8081",
		ErrorLog:     sl,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	l.Info("Starting image server at", "host", server.Addr)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Error("Server stopped working due to error", "err", err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	<-sigChan

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)

}
