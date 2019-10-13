package microsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

// MakeHealthEndpoint - returns an endpoint for the health function
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return models.HealthResponse{Health: s.Health()}, nil
	}
}

// endpoints.txt
