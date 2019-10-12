package microsvc

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

// ServiceDiscovery - returns a service discovery registrar
func (s service) ServiceDiscovery(
	address string,
	port string,
) (
	*api.Client,
	*api.AgentServiceRegistration,
	error,
) {
	u, err := url.Parse(
		fmt.Sprintf(
			"%s:%s",
			s.config.Consul.ConsulAddr,
			s.config.Consul.ConsulPort,
		),
	)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("ADD", address)
	consulAddress := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	client, err := api.NewClient(&api.Config{
		Address: consulAddress,
		Scheme:  u.Scheme,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: s.config.Consul.ConsulToken,
	})
	if err != nil {
		return nil, nil, err
	}
	rand.Seed(time.Now().UTC().UnixNano())
	id := fmt.Sprintf("%s-%d", s.name(), rand.Intn(100))

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, nil, err
	}

	return client, &api.AgentServiceRegistration{
		ID:   id,
		Name: s.name(),
		Tags: []string{
			fmt.Sprintf("service_id=%s", id),
		},
		Port:    p,
		Address: address,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			Name:                           "HTTP API service /health check",
			Method:                         "GET",
			Interval:                       "1s",
			Timeout:                        "60s",
			HTTP:                           fmt.Sprintf("http://%s:%d/health", address, p),
			Notes:                          "service health check",
		},
	}, nil
}
