package grpc

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/php-grpc/tests"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/env"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	ngrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"testing"
	"time"
)

type testCfg struct {
	grpcCfg string
	rpcCfg  string
	envCfg  string
	target  string
}

func (cfg *testCfg) Get(name string) service.Config {
	if name == ID {
		if cfg.grpcCfg == "" {
			return nil
		}

		return &testCfg{target: cfg.grpcCfg}
	}

	if name == rpc.ID {
		return &testCfg{target: cfg.rpcCfg}
	}

	if name == env.ID {
		return &testCfg{target: cfg.envCfg}
	}

	return nil
}

func (cfg *testCfg) Unmarshal(out interface{}) error {
	return json.Unmarshal([]byte(cfg.target), out)
}

func Test_Service_NoConfig(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.Error(t, c.Init(&testCfg{grpcCfg: `{}`}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusRegistered, st)
}

func Test_Service_Configure_Enable(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{grpcCfg: `{
			"listen": "tcp://:9080",
			"tls": {
				"key": "tests/server.key",
				"cert": "tests/server.crt"
			},
			"proto": "tests/test.proto",
			"workers":{
				"command": "php tests/worker.php",
				"relay": "pipes",
				"pool": {
					"numWorkers": 1, 
					"allocateTimeout": 1,
					"destroyTimeout": 1 
				}
			}
	}`}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)
}

func Test_Service_Echo(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{grpcCfg: `{
			"listen": "tcp://:9080",
			"tls": {
				"key": "tests/server.key",
				"cert": "tests/server.crt"
			},
			"proto": "tests/test.proto",
			"workers":{
				"command": "php tests/worker.php",
				"relay": "pipes",
				"pool": {
					"numWorkers": 1, 
					"allocateTimeout": 10,
					"destroyTimeout": 10 
				}
			}
	}`}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := makeClient("localhost:9080")
	defer cn.Close()

	out, err := cl.Echo(context.Background(), &tests.Message{Msg: "ping"})

	assert.NoError(t, err)
	assert.Equal(t, "ping", out.Msg)
}

func makeClient(addr string) (client tests.TestClient, conn *ngrpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile("tests/server.crt", "")
	if err != nil {
		panic(err)
	}

	conn, err = ngrpc.Dial(addr, ngrpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	return tests.NewTestClient(conn), conn
}
