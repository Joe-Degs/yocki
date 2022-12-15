package secrets

import (
	"github.com/Joe-Degs/yocki/internal/config"
	"github.com/Joe-Degs/yocki/server"
	"github.com/hashicorp/vault/api"
)

// Secret interacts with a vault server by binding the api and providing a nice
// interface customized for our needs and uses of vault.
type Secret struct{}

// vaultConfig makes an api.Config out of yocki's vault config
func (s Secret) vaultConfig() *api.Config {
	_ = config.GetConfig()
	return nil
}

func (s *Secret) Routes() []*server.Route {
	return []*server.Route{
		{
			Path:        "",
			Methods:     []string{},
			HandlerFunc: nil,
		},
	}
}

func (Secret) Version() string { return "secret" }
