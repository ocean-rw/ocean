package mongo

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/master/db/common"
)

var _ common.ConfigTableIF = (*ConfigTable)(nil)

const (
	// clusterID 是唯一标示集群的文档 id。
	clusterID = "cluster_id"
)

type ConfigTable struct {
	tbl *mongo.Collection
}

func OpenConfigTable(tbl *mongo.Collection) (*ConfigTable, error) {
	if err := initClusterID(tbl); err != nil {
		return nil, err
	}
	return &ConfigTable{tbl: tbl}, nil
}

func initClusterID(tbl *mongo.Collection) error {
	err := tbl.FindOne(context.Background(), bson.M{"_id": clusterID}).Err()
	if err == mongo.ErrNoDocuments {
		id := uuid.NewString()
		_, err = tbl.InsertOne(context.Background(), bson.M{"_id": clusterID, "uuid": id})
	}
	return err
}

func (t ConfigTable) ClusterID(ctx context.Context) (string, error) {
	var id struct {
		UUID string `bson:"uuid"`
	}
	err := t.tbl.FindOne(ctx, bson.M{"_id": clusterID}).Decode(&id)
	return id.UUID, err
}
