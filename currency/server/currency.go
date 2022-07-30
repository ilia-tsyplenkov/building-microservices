package server

import (
	"context"
	"io"
	"time"

	"github.com/ilia-tsyplenkov/building-microservices/currency/data"
	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
			for i, rr := range v {
				rate, err := c.rates.GetRate(rr.Base.String(), rr.Destination.String())
				if err != nil {
					c.log.Error("unable to get rates", "Base", rr.Base.String(), "Destination", rr.Destination.String())
					// do not serve this RateRequest anymore
					v = append(v[:i], v[i+1:]...)
				}
				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{
							Base: rr.Base, Destination: rr.Destination, Rate: rate},
					},
				})
				if err != nil {
					c.log.Error("unable to send updated rate", "Base", rr.Base.String(), "Destination", rr.Destination.String(), "error", err)
					// delete problem subscription from the map
					delete(c.subscriptions, k)
				}

			}
		}

	}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	if rr.Base == rr.Destination {
		err := status.Newf(codes.InvalidArgument,
			"Base currency - %q can not be the same as the destination one - %q",
			rr.Base.String(),
			rr.Destination.String())
		err, wde := err.WithDetails(rr)
		if wde != nil {
			return nil, wde
		}
		return nil, err.Err()
	}
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
		var validationError *status.Status
		for _, v := range rrs {
			if rr.Base == v.Base && rr.Destination == v.Destination {
				validationError := status.Newf(codes.AlreadyExists,
					"Unable to subscribe for currency as subscription already exists")
				// add original request as metadata
				validationError, err = validationError.WithDetails(rr)
				if err != nil {
					c.log.Error("unable to add metadata to error", "error", err)
				}
				break
			}
		}
		// if a validation error return it
		if validationError != nil {
			src.Send(
				&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_Error{
						Error: validationError.Proto(),
					},
				},
			)
		} else {
			rrs = append(rrs, rr)
			c.subscriptions[src] = rrs
		}
	}

	return nil
}
