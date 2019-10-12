package microsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// LoggingMiddleware - returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if v := ctx.Value("jwt"); v != nil {
					fmt.Println("found value:", v)
				}
				id := ctx.Value(kithttp.ContextKeyRequestXRequestID).(string)
				logger.Log(
					"service", Name,
					"request_id", id,
					"dur", time.Since(begin),
					"mtype", "rate",
					"unit", "µs",
					"error", err,
					"response", response,
					"message", fmt.Sprintf(
						"%s request_id %s event logged with err %v",
						Name, id, err,
					),
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}
