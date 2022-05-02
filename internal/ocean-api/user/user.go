package user

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/pkg/storage"
)

type Config struct {
}

type Mgr struct {
	cfg     *Config
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func New(logger *zap.SugaredLogger, cfg *Config) (*Mgr, error) {
	return &Mgr{logger: logger, cfg: cfg}, nil
}

func (m *Mgr) RegisterRouter(r *chi.Mux) {

}
