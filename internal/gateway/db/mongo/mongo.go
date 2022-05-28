package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
)

type Config struct {
	URI    string `yaml:"uri"`
	DBName string `yaml:"db_name"`

	UserTable   string `yaml:"user_table"`
	BucketTable string `yaml:"bucket_table"`
	FileTable   string `yaml:"file_table"`
}

func Open(cfg *Config) (*common.Database, error) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	protoDB := &common.Database{CloseFn: cli.Disconnect}
	db := cli.Database(cfg.DBName)

	protoDB.UserTable, err = openUserTable(db.Collection(cfg.UserTable))
	if err != nil {
		return nil, err
	}

	protoDB.BucketTable, err = openBucketTable(db.Collection(cfg.BucketTable))
	if err != nil {
		return nil, err
	}

	protoDB.FileTable, err = openFileTable(db.Collection(cfg.FileTable))
	if err != nil {
		return nil, err
	}

	return protoDB, nil
}
