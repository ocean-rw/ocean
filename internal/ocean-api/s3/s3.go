package s3

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/pkg/storage"
)

type Config struct {
	Storage *storage.Config
}

type Service struct {
	cfg     *Config
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func New(logger *zap.SugaredLogger, cfg *Config) (*Service, error) {
	stg, err := storage.New(cfg.Storage)
	if err != nil {
		return nil, err
	}
	return &Service{logger: logger, cfg: cfg, storage: stg}, nil
}

func (s *Service) RegisterRouter(r *chi.Mux) {

}
