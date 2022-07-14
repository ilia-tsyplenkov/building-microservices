package main

import (
	"building-microservices/product-api/data"
	"building-microservices/product-api/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	v := data.NewValidation()
	ph := handlers.NewProducts(l, v)
	r := mux.NewRouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.Use(ph.MiddlewareValidateProduct)
	putRouter.HandleFunc("/products", ph.UpdateProduct)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.Use(ph.MiddlewareValidateProduct)
	postRouter.HandleFunc("/products", ph.Create)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      cors(r),
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
