package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type Server struct {
	address string

	router chi.Router

	srv *http.Server

	logger *zap.Logger
}

func NewServer(address string, logger *zap.Logger, opts ...ServerOption) (s *Server) {
	s = &Server{
		address: address,

		router: chi.NewRouter(),

		logger: logger,
	}

	// fork gin-zap and implement it for chi
	s.router.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer)

	s.srv = &http.Server{
		Addr:              s.address,
		Handler:           s.router,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func (s *Server) ListenAndServe() (err error) {
	s.logger.Info("starting")

	if err = s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Close(ctx context.Context) (err error) {
	return s.srv.Shutdown(ctx)
}
