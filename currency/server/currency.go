package server

import (
	"context"
	"io"
	"time"

	"github.com/ilia-tsyplenkov/building-microservices/currency/data"
	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
	protos.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	currency := &Currency{rates: r, log: l, subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go currency.handleUpdates()
	return currency
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		for k, v := range c.subscriptions {
			for _, rr := range v {
				rate, err := c.rates.GetRate(rr.Base.String(), rr.Destination.String())
				if err != nil {
					c.log.Error("unable to get rates", "Base", rr.Base.String(), "Destination", rr.Destination.String())
				}
				if err := k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}); err != nil {
					c.log.Error("unable to send updated rate", "Base", rr.Base.String(), "Destination", rr.Destination.String())
				}

			}
		}

	}
}
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("client has closed the connection.")
			break
		}
		if err != nil {
			c.log.Error("unable to read from client", "error", err)
			return err
		}
		c.log.Info("handle client request", "Base", rr.Base.String(), "Destionation", rr.Destination.String())

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = make([]*protos.RateRequest, 0)
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
