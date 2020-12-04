package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/morozovcookie/sea-battle/cmd/sea-battle/config"
	"github.com/morozovcookie/sea-battle/http"
	httpV1 "github.com/morozovcookie/sea-battle/http/v1"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	appname = "sea_battle_server"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("init zap logger error", err)
	}

	logger = logger.With(zap.String("appname", appname))

	var s *http.Server
	{
		cfg := config.New()
		if err = cfg.Parse(); err != nil {
			logger.Fatal("parse config error", zap.Error(err))
		}

		seaBattleSvc := httpV1.NewSeaBattleService(
			logger.With(zap.String("component", "sea_battle_http_service")))

		s = http.NewServer(
			cfg.ServerConfig.Address,
			logger.With(zap.String("component", "http_server")),
			http.WithHandler(httpV1.SeaBattleSvcPathPrefix, seaBattleSvc))
	}

	logger.Info("starting application")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(s.ListenAndServe)

	logger.Info("application started")

	select {
	case <-quit:
		break
	case <-ctx.Done():
		break
	}

	logger.Info("stopping application")

	if err = s.Close(context.Background()); err != nil {
		logger.Error("error while closing http server", zap.Error(err))
	}

	if err = eg.Wait(); err != nil {
		logger.Error("error while stopping application", zap.Error(err))
	}

	logger.Info("application stopped")
}
