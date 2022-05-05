package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/lib/config"
	"github.com/ocean-rw/ocean/internal/lib/log"
	"github.com/ocean-rw/ocean/internal/ocean-mgr/db"
	"github.com/ocean-rw/ocean/internal/ocean-mgr/service"
)

type Config struct {
	BindAddr string          `yaml:"bind_addr"`
	Log      *log.Config     `yaml:"log"`
	DB       *db.Config      `yaml:"db"`
	Mgr      *service.Config `yaml:"mgr"`
}

func main() {
	cfgFile := flag.String("c", "config.yaml", "The path to config file.")
	flag.Parse()
	zap.S().Infof("current config file %s", *cfgFile)
	cfg := new(Config)
	err := config.Load(cfg, *cfgFile)
	if err != nil {
		zap.S().Fatalf("failed to load config %s, err: %s", *cfgFile, err)
	}

	logger := log.New(cfg.Log)
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(log.Middleware(logger))

	database, err := db.Open(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect database, err: %s", err)
	}
	defer database.CloseFn(context.TODO())

	s, err := service.New(cfg.Mgr, logger, database)
	if err != nil {
		logger.Fatalf("failed to new disk mgr, err: %s", err)
	}
	s.RegisterRouter(r)

	logger.Infof("ocean-mgr is running at %s", cfg.BindAddr)
	if err = http.ListenAndServe(cfg.BindAddr, r); err != nil {
		logger.Errorf("server closed, err: %s", err)
	}
}
