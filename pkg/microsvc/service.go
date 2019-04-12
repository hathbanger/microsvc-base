package microsvc

import (
	"context"
	"errors"
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

// Service - interface into service methods
type Service interface {
	Health() bool
	ServiceDiscovery() (*api.Client, *api.AgentServiceRegistration, error)

	Foo(context.Context, models.FooRequest) (models.FooResponse, error)
}

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

func (s service) Foo(
	ctx context.Context,
	req models.FooRequest,
) (models.FooResponse, error) {
	var (
		response models.FooResponse
		err      error
	)

	if req.Str == "" {
		err = errors.New("no string was passed")
		s.logger.Log("err", "boo")
		return response, err
	}

	product := req.Str + "bar"
	response.Res = product
	s.logger.Log("response", response.Res)

	return response, err
}
