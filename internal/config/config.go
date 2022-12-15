package config

import "github.com/BurntSushi/toml"

// config for certificate and their management
type Certs struct {
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

// config for the vault mediator service, it contains config for connecting
// to a vault server or agent.
type VaultConfig struct {
	Address      string `toml:"address"`
	AgentAddress string `toml:"agent_address"`

	// specify whether to use custom client, this is usually the case if you
	// want to authenticate with a tls server
	CustomClient bool  `toml:"custom_client"`
	Certs        Certs `toml:"certs"`
}

// Configuration variables for the server and services that it offers
type Config struct {
	Address     string `toml:"listener_address"`
	Certs       Certs  `toml:"certs"`
	VaultConfig `toml:"vault"`
}

var config = &Config{}

func GetConfig() Config { return *config }

// Initialize the config by loading the .toml file specified.
func Init(cpath string) error {
	if _, err := toml.DecodeFile(cpath, config); err != nil {
		return err
	}
	return nil
}
