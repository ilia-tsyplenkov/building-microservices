package main

import (
	"net"
	"os"

	"github.com/ilia-tsyplenkov/building-microservices/currency/data"
	"github.com/ilia-tsyplenkov/building-microservices/currency/server"

	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	gc := grpc.NewServer()
	rates, err := data.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(rates, log)
	protos.RegisterCurrencyServer(gc, cs)
	reflection.Register(gc)
	addr := ":8082"
	log.Info("Starting gRPC server", "address", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gc.Serve(l)
}
