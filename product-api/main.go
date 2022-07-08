package main

import (
	"building-microservices/product-api/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)
	// sm := http.NewServeMux()
	r := mux.NewRouter()
	//sm.Handle("/products", ph).Methods("GET").Subrouter()
	getRouter := r.PathPrefix("/").Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := r.PathPrefix("/").Methods("PUT").Subrouter()
	putRouter.Use(ph.MiddlewareValidateProduct)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	postRouter := r.PathPrefix("/").Methods("POST").Subrouter()
	postRouter.Use(ph.MiddlewareValidateProduct)
	postRouter.HandleFunc("/", ph.AddProduct)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	l.Printf("Starting server at %s ...\n", server.Addr)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}

	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	sig := <-sigChan
	l.Println("Got signal -", sig)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	l.Println("Shutting server down...")
	server.Shutdown(ctx)

}
