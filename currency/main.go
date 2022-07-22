package main

import (
	"net"
	"os"

	"github.com/ilia-tsyplenkov/building-microservices/currency/server"

	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	gc := grpc.NewServer()
	cs := server.NewCurrency(log)
	protos.RegisterCurrencyServer(gc, cs)
	reflection.Register(gc)
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gc.Serve(l)
}
