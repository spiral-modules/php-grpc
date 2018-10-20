package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/php-grpc/tests"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func Test_Proxy_Error(t *testing.T) {
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

	cl, cn := getClient("localhost:9080")
	defer cn.Close()

	_, err := cl.Throw(context.Background(), &tests.Message{Msg: "notFound"})

	assert.Error(t, err)
	se, _ := status.FromError(err)
	assert.Equal(t, "nothing here", se.Message())
	assert.Equal(t, codes.NotFound, se.Code())
}

func Test_Proxy_Metadata(t *testing.T) {
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

	cl, cn := getClient("localhost:9080")
	defer cn.Close()

	out, err := cl.Info(
		metadata.AppendToOutgoingContext(context.Background(), "key", "proxy-value"),
		&tests.Message{Msg: "MD"},
	)

	assert.NoError(t, err)
	assert.Equal(t, `["proxy-value"]`, out.Msg)
}
