package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
	"github.com/ocean-rw/ocean/pkg/proto"
)

var _ common.FileTableIF = (*fileTable)(nil)

type fileTable struct {
	tbl *mongo.Collection
}

func openFileTable(tbl *mongo.Collection) (*fileTable, error) {
	return &fileTable{tbl: tbl}, nil
}

func (t *fileTable) Upsert(ctx context.Context, fileInfo *proto.FileInfo) error {
	fileInfo.PutTime = time.Now()
	filter := bson.M{"_id": fileInfo.FID, "bid": fileInfo.BID}
	_, err := t.tbl.ReplaceOne(ctx, filter, fileInfo, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (t *fileTable) Get(ctx context.Context, bid string, fid string) (*proto.FileInfo, error) {
	filter := bson.M{"_id": fid, "bid": bid}
	ret := new(proto.FileInfo)
	err := t.tbl.FindOne(ctx, filter).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *fileTable) List(ctx context.Context, bid string) ([]*proto.FileInfo, error) {
	filter := bson.M{"bid": bid}
	cur, err := t.tbl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	ret := make([]*proto.FileInfo, 0)
	err = cur.All(ctx, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
