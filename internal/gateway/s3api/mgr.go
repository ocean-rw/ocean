package s3api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
	"github.com/ocean-rw/ocean/pkg/storage"
)

const iso8601TimeFormat = "2006-01-02T15:04:05.000Z"

type Config struct {
	Storage *storage.Config
}

type Mgr struct {
	cfg     *Config
	logger  *zap.SugaredLogger
	db      *common.Database
	storage storage.Storage
}

func New(cfg *Config, logger *zap.SugaredLogger, db *common.Database) (*Mgr, error) {
	stg, err := storage.New(cfg.Storage)
	if err != nil {
		return nil, err
	}
	return &Mgr{
		cfg:     cfg,
		logger:  logger,
		db:      db,
		storage: stg,
	}, nil
}

func (m *Mgr) RegisterRouter(r *chi.Mux) {
	r.Get("/", m.ListBuckets)

	r.Route("/{bucket_id:.+}", func(r chi.Router) {
		r.Put("/", m.CreateBucket)
		r.Delete("/", m.DeleteBucket)
		r.Get("/", m.ListObjects)

		r.Put("/{object_id:.+}", m.PutObject)
		r.Get("/{object_id:.+}", m.GetObject)
		r.Delete("/{object_id:.+}", m.DeleteObject)
	})
}
