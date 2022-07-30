package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
	"github.com/ilia-tsyplenkov/building-microservices/product-api/handlers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	//l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	l := hclog.Default()
	v := data.NewValidation()
	grpcConn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error("error estalishing a connection to grpc server", "error", err)
		panic(err)
	}
	defer grpcConn.Close()
	// create a currency client
	cc := protos.NewCurrencyClient(grpcConn)
	db := data.NewProductDB(cc, l)

	ph := handlers.NewProducts(l, v, db)
	r := mux.NewRouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products", ph.GetProducts).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct).Queries("currency", "{[A-Z]{3}}")

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
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	l.Info("Starting server at", "address", server.Addr)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Error("server finished with an error", "error", err)
		}

	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	sig := <-sigChan
	l.Info("Got signal", "signal", sig)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	l.Info("Shutting server down...")
	server.Shutdown(ctx)
}
