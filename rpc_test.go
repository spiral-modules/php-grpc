package grpc

import (
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/php-grpc/tests"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_RPC(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(rpc.ID, &rpc.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		rpcCfg: `{"enable":true, "listen":"tcp://:5004"}`,
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

	s, _ := c.Get(ID)
	ss := s.(*Service)

	s2, _ := c.Get(rpc.ID)
	rs := s2.(*rpc.Service)

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	rcl, err := rs.Client()
	assert.NoError(t, err)

	out, err := cl.Info(context.Background(), &tests.Message{Msg: "PID"})

	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(*ss.rr.Workers()[0].Pid), out.Msg)

	r := ""
	assert.NoError(t, rcl.Call("grpc.Reset", true, &r))
	assert.Equal(t, "OK", r)

	out2, err := cl.Info(context.Background(), &tests.Message{Msg: "PID"})

	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(*ss.rr.Workers()[0].Pid), out2.Msg)

	assert.NotEqual(t, out.Msg, out2.Msg)
}

func Test_Workers(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(rpc.ID, &rpc.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		rpcCfg: `{"enable":true, "listen":"tcp://:5004"}`,
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

	s, _ := c.Get(ID)
	ss := s.(*Service)

	s2, _ := c.Get(rpc.ID)
	rs := s2.(*rpc.Service)

	go func() { assert.NoError(t, c.Serve()) }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	cl, cn := getClient(addr)
	defer cn.Close()

	rcl, err := rs.Client()
	assert.NoError(t, err)

	out, err := cl.Info(context.Background(), &tests.Message{Msg: "PID"})

	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(*ss.rr.Workers()[0].Pid), out.Msg)

	r := &WorkerList{}
	assert.NoError(t, rcl.Call("grpc.Workers", true, &r))
	assert.Len(t, r.Workers, 1)
}

func Test_Errors(t *testing.T) {
	r := &rpcServer{nil}

	assert.Error(t, r.Reset(true, nil))
	assert.Error(t, r.Workers(true, nil))
}
