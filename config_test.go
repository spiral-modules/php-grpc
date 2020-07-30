package grpc

import (
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
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
			Key:    "tests/server.key",
			Cert:   "tests/server.crt",
			RootCA: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
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

func Test_Config_No_Proto(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test2.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
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

func Test_Config_BadAddress(t *testing.T) {
	cfg := &Config{
		Listen: "tcp//8080",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
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

func Test_Config_BadListener(t *testing.T) {
	cfg := &Config{
		Listen: "unix//8080",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	_, err := cfg.Listener()
	assert.Error(t, err)
}

func Test_Config_UnixListener(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on " + runtime.GOOS)
	}

	cfg := &Config{
		Listen: "unix://rr.sock",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test2.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	ln, err := cfg.Listener()
	assert.NoError(t, err)
	assert.NotNil(t, ln)
	ln.Close()
}

func Test_Config_InvalidWorkerPool(t *testing.T) {
	cfg := &Config{
		Listen: "unix://rr.sock",
		TLS: TLS{
			Key:  "tests/server.key",
			Cert: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				AllocateTimeout: 0,
			},
		},
	}

	assert.Error(t, cfg.Valid())
}

func Test_Config_TLS_No_key(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Cert: "tests/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	// should return nil, because c.EnableTLS will be false in case of missed certs
	assert.Nil(t, cfg.Valid())
}

func Test_Config_TLS_WrongKeyPath(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Cert: "testss/server.crt",
			Key:  "testss/server.key",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
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

func Test_Config_TLS_WrongRootCAPath(t *testing.T) {
	cfg := &Config{
		Listen: "tcp://:8080",
		TLS: TLS{
			Cert:   "tests/server.crt",
			Key:    "tests/server.key",
			RootCA: "testss/server.crt",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
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
			Key: "tests/server.key",
		},
		Proto: "tests/test.proto",
		Workers: &roadrunner.ServerConfig{
			Command: "php tests/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	}

	// should return nil, because c.EnableTLS will be false in case of missed certs
	assert.Nil(t, cfg.Valid())
}
