package microsvc

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

func (s service) ServiceDiscovery() (
	*api.Client,
	*api.AgentServiceRegistration,
	error,
) {

	rand.Seed(time.Now().UTC().UnixNano())

	// err := s.prepareCertificates()
	// if err != nil {
	// 	return nil, nil, err
	// }

	d, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

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

	address := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())

	s.logger.Log(
		"host", u.Hostname(),
		"port", u.Port(),
		"scheme", u.Scheme,
		"address", fmt.Sprintf(
			"%s:%s",
			s.config.Consul.ConsulAddr,
			s.config.Consul.ConsulPort,
		),
		"addr", address,
	)

	tls, err := api.SetupTLSConfig(
		&api.TLSConfig{
			Address:            address,
			CAFile:             fmt.Sprintf("%s/certs/consul_ca.cert", d),
			CAPath:             fmt.Sprintf("%s/certs", d),
			CertFile:           fmt.Sprintf("%s/certs/client.cert", d),
			KeyFile:            fmt.Sprintf("%s/certs/client.key", d),
			InsecureSkipVerify: true,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	c := &api.Config{
		Address: address,
		Scheme:  u.Scheme,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tls,
			},
			Timeout: 30 * time.Second,
		},
		Token: s.config.Consul.ConsulToken,
	}

	client, err := api.NewClient(c)
	if err != nil {
		return nil, nil, err
	}

	port, err := strconv.Atoi(s.config.ServicePort)
	if err != nil {
		return nil, nil, err
	}

	return client, &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%d", s.name(), rand.Intn(100)),
		Name:    fmt.Sprintf("%s", s.name()),
		Address: os.Getenv("CF_INSTANCE_ADDR"),
		Port:    port,
		Tags: []string{
			fmt.Sprintf("cf_instance_addr=%s", os.Getenv("CF_INSTANCE_ADDR")),
		},
		Check: &api.AgentServiceCheck{
			HTTP: fmt.Sprintf(
				"%s/health",
				s.config.ServiceAddr,
			),
			Interval: "10s",
			Timeout:  "1s",
			Notes:    "service health check",
		},
	}, nil
}

// func (s service) prepareCertificates() error {
// 	d, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
// 	s.logger.Log("dir", d)
//
// 	err = s.write(
// 		fmt.Sprintf("%s/certs/client.cert", d),
// 		[]byte(s.config.Consul.ConsulClientCert),
// 	)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = s.write(
// 		fmt.Sprintf("%s/certs/client.key", d),
// 		[]byte(s.config.Consul.ConsulClientKey),
// 	)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = s.write(
// 		fmt.Sprintf("%s/certs/consul_ca.cert", d),
// 		[]byte(s.config.Consul.ConsulCACert),
// 	)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (s service) write(path string, b []byte) error {
	f, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY, 0600,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
