package grpc

import "github.com/spiral/roadrunner/service"

type Config struct {
	// Proto file associated with the service
	Proto string
}

// Hydrate the config and validate it's values.
func (c *Config) Hydrate(cfg service.Config) error {
	return cfg.Unmarshal(c)
}
