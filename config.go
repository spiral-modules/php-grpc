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

// Config describes GRPC service configuration.
type Config struct {
	// Address to listen.
	Listen string

	// Proto file associated with the service.
	Proto string

	// TLS defined authentication method (TLS for now).
	TLS TLS

	// Workers configures roadrunner grpc and worker pool.
	Workers *roadrunner.ServerConfig
}

// TLS defines auth credentials.
type TLS struct {
	// Key defined private server key.
	Key string

	// Cert is https certificate.
	Cert string
}

// Hydrate the config and validate it's values.
func (c *Config) Hydrate(cfg service.Config) error {
	c.Workers = &roadrunner.ServerConfig{}
	c.Workers.InitDefaults()

	if err := cfg.Unmarshal(c); err != nil {
		return err
	}
	c.Workers.UpscaleDurations()

	return c.Valid()
}

// Valid validates the configuration.
func (c *Config) Valid() error {
	if c.Proto == "" && c.Workers.Command != "" {
		// only when rr server is set
		return errors.New("proto file is required")
	}

	if c.Proto != "" {
		if _, err := os.Stat(c.Proto); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("proto file '%s' does not exists", c.Proto)
			}

			return err
		}
	}

	if c.Workers.Command != "" {
		if err := c.Workers.Pool.Valid(); err != nil {
			return err
		}
	}

	if !strings.Contains(c.Listen, ":") {
		return errors.New("mailformed grpc grpc address")
	}

	if c.EnableTLS() {
		if _, err := os.Stat(c.TLS.Key); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("key file '%s' does not exists", c.TLS.Key)
			}

			return err
		}

		if _, err := os.Stat(c.TLS.Cert); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("cert file '%s' does not exists", c.TLS.Cert)
			}

			return err
		}
	}

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

// EnableTLS returns true if rr must listen TLS connections.
func (c *Config) EnableTLS() bool {
	return c.TLS.Key != "" || c.TLS.Cert != ""
}
