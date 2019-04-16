package microsvc

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

const (
	// ServiceName - name of the service
	ServiceName = "microsvc-base"
)

//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 . Service

type service struct {
	mut    *sync.Mutex
	config *models.Config
	logger log.Logger
}

// New - returns new service
func New(config *models.Config, logger log.Logger) Service {
	return service{
		mut:    &sync.Mutex{},
		config: config,
		logger: logger,
	}
}

func (s service) name() string {
	return ServiceName
}

func (s service) Health() bool {
	return true
}

// Service - interface into service methods
type Service interface {
	Health() bool
	ServiceDiscovery() (*api.Client, *api.AgentServiceRegistration, error)

	Foo(context.Context, models.FooRequest) (models.FooResponse, error)
	// here
}
