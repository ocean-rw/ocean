package s3

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/pkg/storage"
)

type Config struct {
	Storage *storage.Config
}

type Mgr struct {
	cfg     *Config
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func New(logger *zap.SugaredLogger, cfg *Config) (*Mgr, error) {
	stg, err := storage.New(cfg.Storage)
	if err != nil {
		return nil, err
	}
	return &Mgr{logger: logger, cfg: cfg, storage: stg}, nil
}

func (m *Mgr) RegisterRouter(r *chi.Mux) {

}
