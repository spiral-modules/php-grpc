package grpc

import (
	"encoding/json"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockCfg struct{ cfg string }

func (cfg *mockCfg) Get(name string) service.Config  { return nil }
func (cfg *mockCfg) Unmarshal(out interface{}) error { return json.Unmarshal([]byte(cfg.cfg), out) }

func Test_Config_Hydrate_Error2(t *testing.T) {
	cfg := &mockCfg{`{"`}
	c := &Config{}

	assert.Error(t, c.Hydrate(cfg))
}

func Test_Config_Valid_TLS(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/server.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	assert.NoError(t, cfg.Valid())
	assert.True(t, cfg.EnableTLS())
}

func Test_Config_TLS_No_key(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Cert: "fixtures/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/server.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	assert.Error(t, cfg.Valid())
}

func Test_Config_TLS_No_Cert(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Key: "fixtures/server.key",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/server.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	assert.Error(t, cfg.Valid())
}
