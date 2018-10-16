package grpc

import "sync"

const ID = "grpc"

type Service struct {
	cgf  *Config
	mu   sync.Mutex
	stop chan interface{}
}

// Init service.
func (s *Service) Init(cfg *Config) (bool, error) {
	s.cgf = cfg

	// register service

	return true, nil
}

// Serve GRPC server.
func (s *Service) Serve() error {
	s.mu.Lock()
	s.stop = make(chan interface{})
	s.mu.Unlock()

	<-s.stop
	return nil
}

// Stop the service.
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stop != nil {
		close(s.stop)
		s.stop = nil
	}
}
