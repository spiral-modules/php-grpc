package grpc

import (
	"fmt"
	"github.com/spiral/grpc/parser"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"sync"
)

const ID = "grpc"

type Service struct {
	cfg  *Config
	env  env.Environment
	list []func(event int, ctx interface{})
	mu   sync.Mutex
	rr   *roadrunner.Server
	grpc *grpc.Server
}

// AddListener attaches grpc event watcher.
func (s *Service) AddListener(l func(event int, ctx interface{})) {
	s.list = append(s.list, l)
}

// Init service.
func (s *Service) Init(cfg *Config, e env.Environment) (ok bool, err error) {
	s.cfg = cfg
	s.env = e

	return true, nil
}

// Serve GRPC grpc.
func (s *Service) Serve() error {
	s.mu.Lock()

	lis, err := s.cfg.Listener()
	if err != nil {
		return err
	}

	if s.env != nil {
		values, err := s.env.GetEnv()
		if err != nil {
			return err
		}

		for k, v := range values {
			s.cfg.Workers.SetEnv(k, v)
		}

		s.cfg.Workers.SetEnv("RR_GRPC", "true")
	}

	s.rr = roadrunner.NewServer(s.cfg.Workers)
	s.rr.Listen(s.throw)

	s.grpc = grpc.NewServer(s.serverOptions()...)

	// register services
	if services, err := parser.File(s.cfg.Proto); err != nil {
		return err
	} else {
		for _, service := range services {
			NewProxy(fmt.Sprintf("%s.%s", service.Package, service.Name), s.cfg.Proto, s.rr).Attach(s.grpc)
		}
	}

	// todo: register external service

	s.mu.Unlock()

	if err := s.rr.Start(); err != nil {
		return err
	}
	defer s.rr.Stop()

	return s.grpc.Serve(lis)
}

// Stop the service.
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.grpc == nil {
		return
	}

	go s.grpc.GracefulStop()
}

// throw handles service, grpc and pool events.
func (s *Service) throw(event int, ctx interface{}) {
	for _, l := range s.list {
		l(event, ctx)
	}

	if event == roadrunner.EventServerFailure {
		// underlying rr grpc is dead
		s.Stop()
	}
}

func (s *Service) serverOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		// wrap default proto codec to bypass message marshal/unmarshal when rr is target
		grpc.CustomCodec(&codec{encoding.GetCodec("proto")}),
	}
}
