package grpc

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service"
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

	// see .rr.yaml for the explanations
	MaxSendMsgSize        int64
	MaxRecvMsgSize        int64
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	MaxConcurrentStreams  int64
	PingTime              time.Duration
	Timeout               time.Duration
}

// TLS defines auth credentials.
type TLS struct {
	// Key defined private server key.
	Key string

	// Cert is https certificate.
	Cert string

	// Root CA
	RootCA string
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

		// RootCA is optional, but if provided - check it
		if c.TLS.RootCA != "" {
			if _, err := os.Stat(c.TLS.RootCA); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("root ca path provided, but key file '%s' does not exists", c.TLS.RootCA)
				}
				return err
			}
		}
	}

	// used to set max time
	infinity := time.Duration(math.MaxInt64)

	if c.PingTime == 0 {
		c.PingTime = time.Hour * 2
	}

	if c.Timeout == 0 {
		c.Timeout = time.Second * 20
	}

	if c.MaxConcurrentStreams == 0 {
		c.MaxConcurrentStreams = 10
	}
	// set default
	if c.MaxConnectionAge == 0 {
		c.MaxConnectionAge = infinity
	}

	// set default
	if c.MaxConnectionIdle == 0 {
		c.MaxConnectionIdle = infinity
	}

	if c.MaxConnectionAgeGrace == 0 {
		c.MaxConnectionAgeGrace = infinity
	}

	if c.MaxRecvMsgSize == 0 {
		c.MaxRecvMsgSize = 1024 * 1024 * 50
	} else {
		c.MaxRecvMsgSize = 1024 * 1024 * c.MaxRecvMsgSize
	}

	if c.MaxSendMsgSize == 0 {
		c.MaxSendMsgSize = 1024 * 1024 * 50
	} else {
		c.MaxSendMsgSize = 1024 * 1024 * c.MaxSendMsgSize
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
	// Key and Cert OR Key and Cert and RootCA
	return (c.TLS.RootCA != "" && c.TLS.Key != "" && c.TLS.Cert != "") || (c.TLS.Key != "" && c.TLS.Cert != "")
}
