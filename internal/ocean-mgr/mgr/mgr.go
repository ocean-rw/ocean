package mgr

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
)

type Config struct {
}

type Mgr struct {
	cfg       *Config
	logger    *zap.SugaredLogger
	db        *common.Database
	clusterID string
}

func New(cfg *Config, logger *zap.SugaredLogger, db *common.Database) (*Mgr, error) {
	clusterID, err := db.ConfigTable.ClusterID(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Mgr{
		cfg:       cfg,
		logger:    logger,
		db:        db,
		clusterID: clusterID,
	}, nil
}

func (m *Mgr) RegisterRouter(r *chi.Mux) {
	r.Post("/allocstripes", m.AllocStripes)

	r.Post("/allocdisklabel", m.AllocDiskLabel)
	r.Post("/registerdisks", m.RegisterDisks)
	r.Get("/disks", m.ListDisks)
	r.Get("/disk/{disk_id}", m.GetDisk)
}
