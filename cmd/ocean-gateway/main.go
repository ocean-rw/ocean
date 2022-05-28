package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/gateway/db"
	"github.com/ocean-rw/ocean/internal/gateway/s3api"
	"github.com/ocean-rw/ocean/internal/gateway/user"
	"github.com/ocean-rw/ocean/internal/lib/config"
	"github.com/ocean-rw/ocean/internal/lib/log"
)

type Config struct {
	BindAddr string        `yaml:"bind_addr"`
	Log      *log.Config   `yaml:"log"`
	DB       *db.Config    `yaml:"db"`
	Gateway  *s3api.Config `yaml:"gateway"`
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
	defer func() { _ = logger.Sync() }()

	r := chi.NewRouter()

	database, err := db.Open(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect database, err: %s", err)
	}
	defer func() { _ = database.CloseFn(context.TODO()) }()

	userMgr, err := user.New(database.UserTable)
	if err != nil {
		logger.Fatalf("failed to new api service, err: %s", err)
	}

	s3Mgr, err := s3api.New(cfg.Gateway, logger, database)
	if err != nil {
		logger.Fatalf("failed to new api service, err: %s", err)
	}

	r.Use(log.Middleware(logger), userMgr.Auth())
	userMgr.RegisterRouter(r)
	s3Mgr.RegisterRouter(r)

	logger.Infof("ocean-gateway is running at %s", cfg.BindAddr)
	if err = http.ListenAndServe(cfg.BindAddr, r); err != nil {
		logger.Errorf("server closed, err: %s", err)
	}
}
