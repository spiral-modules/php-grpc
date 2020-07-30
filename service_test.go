package grpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/php-grpc/tests"
	"github.com/spiral/php-grpc/tests/ext"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/env"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	ngrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type testCfg struct {
	grpcCfg string
	rpcCfg  string
	envCfg  string
	target  string
}

func (cfg *testCfg) Get(service string) service.Config {
	if service == ID {
		if cfg.grpcCfg == "" {
			return nil
		}

		return &testCfg{target: cfg.grpcCfg}
	}

	if service == rpc.ID {
		return &testCfg{target: cfg.rpcCfg}
	}

	if service == env.ID {
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
	assert.Equal(t, service.StatusInactive, st)
}

func Test_Service_Configure_Enable(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
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
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)
}

func Test_Service_Dead(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
			"listen": "tcp://:9080",
			"tls": {
				"key": "tests/server.key",
				"cert": "tests/server.crt"
			},
			"proto": "tests/test.proto",
			"workers":{
				"command": "php tests/worker-bad.php",
				"relay": "pipes",
				"pool": {
					"numWorkers": 1, 
					"allocateTimeout": 10,
					"destroyTimeout": 10 
				}
			}
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	assert.Error(t, c.Serve())
}

func Test_Service_Invalid_TLS(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
			"listen": "tcp://:9080",
			"tls": {
				"key": "tests/server.key",
				"cert": "tests/test.proto"
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
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	assert.Error(t, c.Serve())
}

func Test_Service_Invalid_Proto(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
			"listen": "tcp://:9080",
			"tls": {
				"key": "tests/server.key",
				"cert": "tests/server.crt"
			},
			"proto": "tests/server.key",
			"workers":{
				"command": "php tests/worker.php",
				"relay": "pipes",
				"pool": {
					"numWorkers": 1, 
					"allocateTimeout": 10,
					"destroyTimeout": 10 
				}
			}
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	assert.Error(t, c.Serve())
}

func Test_Service_Echo(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
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
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	out, err := cl.Echo(context.Background(), &tests.Message{Msg: "ping"})

	assert.NoError(t, err)
	assert.Equal(t, "ping", out.Msg)
}

func Test_Service_Empty(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
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
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	_, err := cl.Ping(context.Background(), &tests.EmptyMessage{})

	assert.NoError(t, err)
}

func Test_Service_ErrorBuffer(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
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
	}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	goterr := make(chan interface{})
	s.(*Service).AddListener(func(event int, ctx interface{}) {
		if event == roadrunner.EventStderrOutput {
			if string(ctx.([]byte)) == "WORLD\n" {
				goterr <- nil
			}
		}
	})

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	out, err := cl.Die(context.Background(), &tests.Message{Msg: "WORLD"})

	<-goterr
	assert.NoError(t, err)
	assert.Equal(t, "WORLD", out.Msg)
}

func Test_Service_Env(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		grpcCfg: `{
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
	}`,
		envCfg: `
{
	"env_key": "value"
}`,
	}))

	s, st := c.Get(ID)
	assert.NotNil(t, s)
	assert.Equal(t, service.StatusOK, st)

	// should do nothing
	s.(*Service).Stop()

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	out, err := cl.Info(context.Background(), &tests.Message{Msg: "RR_GRPC"})

	assert.NoError(t, err)
	assert.Equal(t, "true", out.Msg)

	out, err = cl.Info(context.Background(), &tests.Message{Msg: "ENV_KEY"})

	assert.NoError(t, err)
	assert.Equal(t, "value", out.Msg)
}

func Test_Service_External_Service_Test(t *testing.T) {
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
	s.(*Service).AddService(func(server *ngrpc.Server) {
		ext.RegisterExternalServer(server, &externalService{})
	})

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getExternalClient("localhost:9080")
	defer cn.Close()

	out, err := cl.Echo(context.Background(), &ext.Ping{Value: 9})

	assert.NoError(t, err)
	assert.Equal(t, int64(90), out.Value)
}

func Test_Service_Kill(t *testing.T) {
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

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)

	s.(*Service).throw(roadrunner.EventServerFailure, nil)
}

func getClient(addr string) (client tests.TestClient, conn *ngrpc.ClientConn) {
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

func getExternalClient(addr string) (client ext.ExternalClient, conn *ngrpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile("tests/server.crt", "")
	if err != nil {
		panic(err)
	}

	conn, err = ngrpc.Dial(addr, ngrpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	return ext.NewExternalClient(conn), conn
}

// externalService service.
type externalService struct{}

// Echo for external service.
func (s *externalService) Echo(ctx context.Context, ping *ext.Ping) (*ext.Pong, error) {
	return &ext.Pong{Value: ping.Value * 10}, nil
}
