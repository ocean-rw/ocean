package user

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/pkg/storage"
)

type Config struct {
}

type Service struct {
	cfg     *Config
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func New(logger *zap.SugaredLogger, cfg *Config) (*Service, error) {
	return &Service{logger: logger, cfg: cfg}, nil
}

func (s *Service) RegisterRouter(r *chi.Mux) {

}
