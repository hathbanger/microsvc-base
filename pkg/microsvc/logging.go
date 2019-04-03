package microsvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// LoggingMiddleware - middeware for logging
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log(
					"service", "microsvc-base",
					"dur", time.Since(begin),
					"mtype", "rate",
					"unit", "s",
					"error", err,
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}
