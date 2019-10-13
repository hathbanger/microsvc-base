package microsvc

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"

	"errors"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

const (
	// ServiceName - name of the service
	ServiceName = "microsvc-base"
)

var (
	// ErrMarshal - error for UnMarshalling
	ErrMarshal = errors.New("could not marshal request")
	// ErrRequest - error if a request cannot be created
	ErrRequest = errors.New("could not create request")
	// ErrToken - error if a token is not present or valid
	ErrToken = errors.New("token is invalid or empty")
	// ErrSize - error for provisioning size configurations
	ErrSize = errors.New("requested provisioning size not found")
	// ErrDNS - error for DNS lookups
	ErrDNS = errors.New("hostnames in use")

	// Arch - the build arch
	Arch string
	// APIVersion - the api version
	APIVersion string
	// BuildTime - the build time
	BuildTime string
	// GitCommit - the git commit
	GitCommit string
	// Name - the service name
	Name = "microsvc-base"
	// Ver - the service version
	Ver string
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
	ServiceDiscovery(string, string) (*api.Client, *api.AgentServiceRegistration, error)

	// interfaceDeclaration.txt
}
