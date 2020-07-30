package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"time"

	"github.com/spiral/php-grpc/parser"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service/env"
	"github.com/spiral/roadrunner/service/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/keepalive"
)

// ID sets public GRPC service ID for roadrunner.Container.
const ID = "grpc"

var errCouldNotAppendPemError = errors.New("could not append Certs from PEM")

// Service manages set of GPRC services, options and connections.
type Service struct {
	cfg      *Config
	env      env.Environment
	list     []func(event int, ctx interface{})
	opts     []grpc.ServerOption
	services []func(server *grpc.Server)
	mu       sync.Mutex
	rr       *roadrunner.Server
	cr       roadrunner.Controller
	grpc     *grpc.Server
}

// Attach attaches cr. Currently only one cr is supported.
func (svc *Service) Attach(ctr roadrunner.Controller) {
	svc.cr = ctr
}

// AddListener attaches grpc event watcher.
func (svc *Service) AddListener(l func(event int, ctx interface{})) {
	svc.list = append(svc.list, l)
}

// AddService would be invoked after GRPC service creation.
func (svc *Service) AddService(r func(server *grpc.Server)) error {
	svc.services = append(svc.services, r)
	return nil
}

// AddOption adds new GRPC server option. Codec and TLS options are controlled by service internally.
func (svc *Service) AddOption(opt grpc.ServerOption) {
	svc.opts = append(svc.opts, opt)
}

// Init service.
func (svc *Service) Init(cfg *Config, r *rpc.Service, e env.Environment) (ok bool, err error) {
	svc.cfg = cfg
	svc.env = e

	if r != nil {
		if err := r.Register(ID, &rpcServer{svc}); err != nil {
			return false, err
		}
	}

	if svc.cfg.Workers.Command != "" {
		svc.rr = roadrunner.NewServer(svc.cfg.Workers)
	}

	return true, nil
}

// Serve GRPC grpc.
func (svc *Service) Serve() (err error) {
	svc.mu.Lock()

	if svc.grpc, err = svc.createGPRCServer(); err != nil {
		svc.mu.Unlock()
		return err
	}

	ls, err := svc.cfg.Listener()
	if err != nil {
		svc.mu.Unlock()
		return err
	}
	defer ls.Close()

	if svc.rr != nil {
		if svc.env != nil {
			if err := svc.env.Copy(svc.cfg.Workers); err != nil {
				svc.mu.Unlock()
				return err
			}
		}

		svc.cfg.Workers.SetEnv("RR_GRPC", "true")

		svc.rr.Listen(svc.throw)

		if svc.cr != nil {
			svc.rr.Attach(svc.cr)
		}

		if err := svc.rr.Start(); err != nil {
			svc.mu.Unlock()
			return err
		}
		defer svc.rr.Stop()
	}

	svc.mu.Unlock()

	return svc.grpc.Serve(ls)
}

// Stop the service.
func (svc *Service) Stop() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	if svc.grpc == nil {
		return
	}

	go svc.grpc.GracefulStop()
}

// Server returns associated rr server (if any).
func (svc *Service) Server() *roadrunner.Server {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	return svc.rr
}

// call info
func (svc *Service) interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)

	svc.throw(EventUnaryCall, &UnaryCallEvent{
		Info:    info,
		Context: ctx,
		Error:   err,
		start:   start,
		elapsed: time.Since(start),
	})

	return resp, err
}

// throw handles service, grpc and pool events.
func (svc *Service) throw(event int, ctx interface{}) {
	for _, l := range svc.list {
		l(event, ctx)
	}

	if event == roadrunner.EventServerFailure {
		// underlying rr grpc is dead
		svc.Stop()
	}
}

// new configured GRPC server
func (svc *Service) createGPRCServer() (*grpc.Server, error) {
	opts, err := svc.serverOptions()
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(opts...)

	if svc.cfg.Proto != "" && svc.rr != nil {
		// php proxy services
		services, err := parser.File(svc.cfg.Proto, path.Dir(svc.cfg.Proto))
		if err != nil {
			return nil, err
		}

		for _, service := range services {
			p := NewProxy(fmt.Sprintf("%s.%s", service.Package, service.Name), svc.cfg.Proto, svc.rr)
			for _, m := range service.Methods {
				p.RegisterMethod(m.Name)
			}

			server.RegisterService(p.ServiceDesc(), p)
		}
	}

	// external and native  services
	for _, r := range svc.services {
		r(server)
	}

	return server, nil
}

// server options
func (svc *Service) serverOptions() (opts []grpc.ServerOption, err error) {
	var tcreds credentials.TransportCredentials
	if svc.cfg.EnableTLS() {
		// if client CA is not empty we combine it with Cert and Key
		if svc.cfg.TLS.RootCA != "" {
			cert, err := tls.LoadX509KeyPair(svc.cfg.TLS.Cert, svc.cfg.TLS.Key)
			if err != nil {
				return nil, err
			}

			certPool, err := x509.SystemCertPool()
			if err != nil {
				return nil, err
			}
			if certPool == nil {
				certPool = x509.NewCertPool()
			}
			rca, err := ioutil.ReadFile(svc.cfg.TLS.RootCA)
			if err != nil {
				return nil, err
			}

			if ok := certPool.AppendCertsFromPEM(rca); !ok {
				return nil, errCouldNotAppendPemError
			}

			tcreds = credentials.NewTLS(&tls.Config{
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{cert},
				ClientCAs:    certPool,
			})
		} else {
			tcreds, err = credentials.NewServerTLSFromFile(svc.cfg.TLS.Cert, svc.cfg.TLS.Key)
			if err != nil {
				return nil, err
			}
		}

		serverOptions := []grpc.ServerOption{
			grpc.MaxSendMsgSize(int(svc.cfg.MaxSendMsgSize)),
			grpc.MaxRecvMsgSize(int(svc.cfg.MaxRecvMsgSize)),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     svc.cfg.MaxConnectionIdle,
				MaxConnectionAge:      svc.cfg.MaxConnectionAge,
				MaxConnectionAgeGrace: svc.cfg.MaxConnectionAge,
				Time:                  svc.cfg.PingTime,
				Timeout:               svc.cfg.Timeout,
			}),
			grpc.MaxConcurrentStreams(uint32(svc.cfg.MaxConcurrentStreams)),
		}

		opts = append(opts, grpc.Creds(tcreds))
		opts = append(opts, serverOptions...)
	}

	opts = append(opts, svc.opts...)

	// custom codec is required to bypass protobuf, common interceptor used for debug and stats
	return append(
		opts,
		grpc.UnaryInterceptor(svc.interceptor),
		grpc.CustomCodec(&codec{encoding.GetCodec("proto")}),
	), nil
}
