package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
)

type Config struct {
	URI    string `yaml:"uri"`
	DBName string `yaml:"db_name"`

	ConfigTable string `yaml:"config_table"`
	IDTable     string `yaml:"id_table"`
	DiskTable   string `yaml:"disk_table"`
	StripeTable string `yaml:"stripe_table"`
}

func Open(cfg *Config) (*common.Database, error) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	protoDB := &common.Database{CloseFn: cli.Disconnect}
	db := cli.Database(cfg.DBName)

	protoDB.ConfigTable, err = OpenConfigTable(db.Collection(cfg.ConfigTable))
	if err != nil {
		return nil, err
	}

	protoDB.IDTable, err = OpenIDTable(db.Collection(cfg.IDTable))
	if err != nil {
		return nil, err
	}

	protoDB.DiskTable, err = OpenDiskTable(db.Collection(cfg.DiskTable))
	if err != nil {
		return nil, err
	}

	protoDB.StripeTable, err = OpenStripeTable(db.Collection(cfg.StripeTable))
	if err != nil {
		return nil, err
	}

	return protoDB, nil
}
