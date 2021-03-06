package main

import (
	"context"
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/lib/config"
	"github.com/ocean-rw/ocean/internal/lib/log"
	"github.com/ocean-rw/ocean/internal/master/db"
	"github.com/ocean-rw/ocean/internal/master/service"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	BindAddr string          `yaml:"bind_addr"`
	Log      *log.Config     `yaml:"log"`
	DB       *db.Config      `yaml:"db"`
	Master   *service.Config `yaml:"master"`
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
	r.Use(log.Middleware(logger))

	database, err := db.Open(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect database, err: %s", err)
	}
	defer func() { _ = database.CloseFn(context.TODO()) }()

	s, err := service.New(cfg.Master, logger, database)
	if err != nil {
		logger.Fatalf("failed to new service, err: %s", err)
	}
	s.RegisterRouter(r)

	logger.Infof("ocean-master is running at %s", cfg.BindAddr)
	if err = http.ListenAndServe(cfg.BindAddr, r); err != nil {
		logger.Errorf("server closed, err: %s", err)
	}
}
