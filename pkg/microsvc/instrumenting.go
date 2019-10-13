package microsvc

import (
	"github.com/go-kit/kit/metrics"
	"github.com/hashicorp/consul/api"
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
func (i instrumentingMiddleware) ServiceDiscovery(address string, port string) (
	*api.Client,
	*api.AgentServiceRegistration,
	error,
) {
	return i.next.ServiceDiscovery(address, port)
}

// instrumenting.txt
