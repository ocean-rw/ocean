package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
	"github.com/ocean-rw/ocean/pkg/proto"
)

var _ common.BucketTableIF = (*bucketTable)(nil)

type bucketTable struct {
	tbl *mongo.Collection
}

func openBucketTable(tbl *mongo.Collection) (*bucketTable, error) {
	return &bucketTable{tbl: tbl}, nil
}

func (t *bucketTable) List(ctx context.Context, uid uint64) ([]*proto.BucketInfo, error) {
	filter := bson.M{"uid": uid}
	cur, err := t.tbl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	ret := make([]*proto.BucketInfo, 0)
	err = cur.All(ctx, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *bucketTable) Insert(ctx context.Context, bucketInfo *proto.BucketInfo) error {
	bucketInfo.PutTime = time.Now()
	_, err := t.tbl.InsertOne(ctx, bucketInfo)
	return err
}

func (t *bucketTable) Get(ctx context.Context, bid string) (*proto.BucketInfo, error) {
	filter := bson.M{"_id": bid}
	ret := new(proto.BucketInfo)
	err := t.tbl.FindOne(ctx, filter).Decode(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *bucketTable) Delete(ctx context.Context, bid string) error {
	filter := bson.M{"_id": bid}
	_, err := t.tbl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
