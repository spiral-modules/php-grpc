package grpc

import (
	"errors"
	"fmt"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service"
	"net"
	"os"
	"strings"
	"syscall"
)

type Config struct {
	// Address to listen.
	Listen string

	// Proto file associated with the service.
	Proto string

	// Auth defined authentication method (TLS for now),
	Auth Auth

	// Workers configures roadrunner grpc and worker pool.
	Workers *roadrunner.ServerConfig
}

// Auth defines auth credentials.
type Auth struct {
}

// Hydrate the config and validate it's values.
func (c *Config) Hydrate(cfg service.Config) error {
	if c.Workers == nil {
		c.Workers = &roadrunner.ServerConfig{}
	}

	c.Workers.InitDefaults()
	if err := cfg.Unmarshal(c); err != nil {
		return err
	}
	c.Workers.UpscaleDurations()

	return c.Valid()
}

// Valid validates the configuration.
func (c *Config) Valid() error {
	if c.Proto == "" {
		return errors.New("proto file is required")
	}

	if _, err := os.Stat(c.Proto); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("proto file '%s' does not exists", c.Proto)
		}

		return err
	}

	if err := c.Workers.Pool.Valid(); err != nil {
		return err
	}

	if !strings.Contains(c.Listen, ":") {
		return errors.New("mailformed grpc grpc address")
	}

	//if c.EnableTLS() {
	//	if _, err := os.Stat(c.SSL.Key); err != nil {
	//		if os.IsNotExist(err) {
	//			return fmt.Errorf("key file '%s' does not exists", c.SSL.Key)
	//		}
	//
	//		return err
	//	}
	//
	//	if _, err := os.Stat(c.SSL.Cert); err != nil {
	//		if os.IsNotExist(err) {
	//			return fmt.Errorf("cert file '%s' does not exists", c.SSL.Cert)
	//		}
	//
	//		return err
	//	}
	//}

	return nil
}

// Listener creates new rpc socket Listener.
func (c *Config) Listener() (net.Listener, error) {
	dsn := strings.Split(c.Listen, "://")
	if len(dsn) != 2 {
		return nil, errors.New("invalid socket DSN (tcp://:6001, unix://rpc.sock)")
	}

	if dsn[0] == "unix" {
		syscall.Unlink(dsn[1])
	}

	return net.Listen(dsn[0], dsn[1])
}
