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

// MakeFooEndpoint - returns endpoint for foo
func MakeFooEndpoint(
	s Service,
	l log.Logger,
	c *models.Config,
) endpoint.Endpoint {
	var e endpoint.Endpoint
	{
		e = func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(models.FooRequest)
			if !ok {
				return nil, ErrBadRequest
			}
			return s.Foo(ctx, req)
		}
		//e = AuthMiddleware([]string{"test"}, "test", c.PublicKey)(e)
		e = LoggingMiddleware(log.With(l, "method", "foo"))(e)
	}
	return e
}

// endpoints.txt
