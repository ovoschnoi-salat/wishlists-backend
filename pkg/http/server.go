package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT"`
	ListenAddress   string `env:"LISTEN_ADDR,required"`
}

type Server struct {
	cfg       ServerConfig
	srv       http.Server
	stopFuncs []func()
}

func NewServer(cfg ServerConfig, h http.Handler) *Server {
	return &Server{cfg: cfg, srv: http.Server{Addr: cfg.ListenAddress, Handler: h}}
}

func (s *Server) Run() error {
	log.Info().Msg("starting server, listening on " + s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) AddStopFunc(f ...func()) {
	s.stopFuncs = append(s.stopFuncs, f...)
}

func (s *Server) Stop() {
	shutdownTimeout := 10 * time.Second
	if s.cfg.ShutdownTimeout > 0 {
		shutdownTimeout = time.Duration(s.cfg.ShutdownTimeout) * time.Second
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := s.srv.Shutdown(shutdownCtx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("error shutting down http server")
	}
	for i := len(s.stopFuncs) - 1; i >= 0; i-- {
		s.stopFuncs[i]()
	}
}
