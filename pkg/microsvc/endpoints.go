package microsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

// MakeHealthEndpoint - returns an endpoint for the health function
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return models.HealthResponse{Health: s.Health()}, nil
	}
}

// MakeFooEndpoint - returns endpoint to foo function
func MakeFooEndpoint(
	s Service,
	l log.Logger,
	mw []endpoint.Middleware,
) endpoint.Endpoint {

	var e endpoint.Endpoint
	e = makeFooEndpoint(s)
	// add middlewares
	for _, m := range mw {
		e = m(e)
	}
	e = LoggingMiddleware(log.With(l, "method", "foo"))(e)
	return e
}

func makeFooEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(models.FooRequest)
		if !ok {
			return nil, ErrBadRequest
		}
		return s.Foo(ctx, req)
	}
}
