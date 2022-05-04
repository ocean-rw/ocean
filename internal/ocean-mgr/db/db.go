package db

import (
	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/mongo"
)

type Config struct {
	Mongo *mongo.Config `yaml:"mongo"`
}

// Open connect to the database,
// and hiding the database implementation details
// for support of multiple databases.
func Open(cfg *Config) (*common.Database, error) {
	if cfg.Mongo != nil {
		return mongo.Open(cfg.Mongo)
	}
	return nil, nil
}
