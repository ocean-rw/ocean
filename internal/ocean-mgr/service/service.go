package service

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
)

type Config struct {
}

type Service struct {
	cfg       *Config
	logger    *zap.SugaredLogger
	db        *common.Database
	clusterID string
}

func New(cfg *Config, logger *zap.SugaredLogger, db *common.Database) (*Service, error) {
	clusterID, err := db.ConfigTable.ClusterID(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Service{
		cfg:       cfg,
		logger:    logger,
		db:        db,
		clusterID: clusterID,
	}, nil
}

func (s *Service) RegisterRouter(r *chi.Mux) {
	r.Post("/allocstripes", s.AllocStripes)

	r.Post("/allocdisklabel", s.AllocDiskLabel)
	r.Post("/registerdisks", s.RegisterDisks)
	r.Get("/disks", s.ListDisks)
	r.Get("/disk/{disk_id}", s.GetDisk)
}
