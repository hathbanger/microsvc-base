package microsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/hashicorp/consul/api"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

// InstrumentingMiddleware - middleware for metrics
func InstrumentingMiddleware(
	duration metrics.Histogram,
	count metrics.Counter,
	service Service,
) Service {
	return &instrumentingMiddleware{
		duration: duration,
		count:    count,
		next:     service,
	}
}

type instrumentingMiddleware struct {
	duration metrics.Histogram
	count    metrics.Counter
	next     Service
}

// Health - instrumentation for the health endpoint
func (i instrumentingMiddleware) Health() bool {
	return i.next.Health()
}

// ServiceDiscovery - service discovery
func (i instrumentingMiddleware) ServiceDiscovery() (
	*api.Client,
	*api.AgentServiceRegistration,
	error,
) {
	return i.next.ServiceDiscovery()
}

// Foo - intrumentation for foo endpoint
func (i instrumentingMiddleware) Foo(
	ctx context.Context,
	request models.FooRequest,
) (res models.FooResponse, err error) {
	defer func(begin time.Time) {
		i.duration.With(
			"method", "Foo",
			"result", fmt.Sprint(err == nil),
			"mtype", "rate",
			"unit", "s",
		).Observe(time.Since(begin).Seconds())
		i.count.With(
			"method", "Foo",
			"result", fmt.Sprint(err == nil),
			"mtype", "count",
			"unit", "req",
		).Add(1)
	}(time.Now())
	return i.next.Foo(ctx, request)
}
