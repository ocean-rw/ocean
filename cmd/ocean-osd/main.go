package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/lib/config"
	"github.com/ocean-rw/ocean/internal/lib/log"
	osd "github.com/ocean-rw/ocean/internal/ocean-osd/mgr"
)

type Config struct {
	BindAddr string      `yaml:"bind_addr"`
	Log      *log.Config `yaml:"log"`
	OSD      *osd.Config `yaml:"osd"`
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

	m, err := osd.New(cfg.OSD, logger)
	if err != nil {
		logger.Fatalf("failed to new disk mgr, err: %s", err)
	}
	m.RegisterRouters(r)

	logger.Infof("ocean-osd is running at %s", cfg.BindAddr)
	if err = http.ListenAndServe(cfg.BindAddr, r); err != nil {
		logger.Fatalf("server closed, err: %s", err)
	}
}
