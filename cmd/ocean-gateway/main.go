package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/gateway/service"
	"github.com/ocean-rw/ocean/internal/lib/config"
	"github.com/ocean-rw/ocean/internal/lib/log"
)

type Config struct {
	BindAddr string          `yaml:"bind_addr"`
	Log      *log.Config     `yaml:"log"`
	Gateway  *service.Config `yaml:"gateway"`
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

	s, err := service.New(logger, cfg.Gateway)
	if err != nil {
		logger.Fatalf("failed to new api service, err: %s", err)
	}
	s.RegisterRouter(r)

	logger.Infof("ocean-gateway is running at %s", cfg.BindAddr)
	if err = http.ListenAndServe(cfg.BindAddr, r); err != nil {
		logger.Errorf("server closed, err: %s", err)
	}
}
