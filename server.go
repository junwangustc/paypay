package main

import log "github.com/junwangustc/ustclog"

type Service interface {
	Open() error
	Close() error
}
type Server struct {
	services []Service
	cfg      *Config
}

func (s *Server) Open() error {
	if err := func() error {
		for _, service := range s.services {
			log.Println("start Opening service %T", service)
			if err := service.Open(); err != nil {
				return err
			}
			log.Println("opened Service %T:%s", service)
		}
		return nil
	}(); err != nil {
		s.Close()
		return err
	}
	return nil
}

func (s *Server) Close() {
	for _, service := range s.services {
		log.Printf("D! closing service: %T", service)
		err := service.Close()
		if err != nil {
			log.Printf("E! error closing service %T: %v", service, err)
		}
		log.Printf("D! closed service: %T", service)
	}
}

func NewServer(cfg *Config) (srv *Server) {
	srv = &Server{cfg: cfg}
	return srv
}
func (s *Server) Run() (err error) {
	if err = s.Open(); err != nil {
		log.Println("[Error]open server fail", err)
		return err
	}
	return nil
}
