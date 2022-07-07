package main

import (
	"building-microservices/product-api/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
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
