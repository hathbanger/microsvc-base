package models

import "crypto/rsa"

// Config - configuration file structure
type Config struct {
	ServiceAddr            string         `json:"service_addr"`
	ServicePort            string         `json:"service_port"`
	HTTPServerReadTimeout  string         `json:"https_server_read_timeout"`
	HTTPServerWriteTimeout string         `json:"https_server_write_timeout"`
	LogPath                string         `json:"log_path,omitempty"`
	PublicKey              *rsa.PublicKey `json:"public_key"`
	Consul                 Consul         `json:"consul"`
	Auth                   Auth           `json:"auth"`
}

// Consul - holds consul config stuff
type Consul struct {
	ConsulAddr       string `json:"consul_addr"`
	ConsulPort       string `json:"consul_port"`
	ConsulToken      string `json:"consul_token"`
	ConsulCACert     string `json:"consul_ca_cert"`
	ConsulClientCert string `json:"consul_client_cert"`
	ConsulClientKey  string `json:"consul_client_key"`
}

// Auth - struc with auth stuff
type Auth struct {
	Groups     []string `json:"groups"`
	ProfileURL string   `json:"profile_url"`
	PublicKey  string   `json:"public_key"`
}
