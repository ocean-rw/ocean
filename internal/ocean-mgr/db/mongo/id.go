package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
)

var _ common.IDTableIF = (*IDTable)(nil)

const (
	// startDiskID 不能为 0，系统认为 diskID = 0 表示空块。
	startDiskID = uint32(10000)
	diskIDTblID = "disk_id"
)

type IDTable struct {
	tbl *mongo.Collection
}

type diskID struct {
	DiskID uint32 `bson:"disk_id"`
}

func OpenIDTable(tbl *mongo.Collection) (*IDTable, error) {
	if err := initDiskIDTable(tbl); err != nil {
		return nil, err
	}
	return &IDTable{tbl: tbl}, nil
}

func initDiskIDTable(tbl *mongo.Collection) error {
	ctx := context.Background()
	var results []bson.M
	err := findAll(ctx, tbl, &results)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		_, err = tbl.InsertOne(ctx, bson.M{"_id": diskIDTblID, "disk_id": startDiskID})
		if mongo.IsDuplicateKeyError(err) {
			err = nil
		}
		return err
	}

	if len(results) != 1 {
		return fmt.Errorf("unexpected diskIDAllocTable size %d", len(results))
	}

	return nil
}

func (t IDTable) AllocDiskID(ctx context.Context) (uint32, error) {
	opt := options.FindOneAndUpdate().SetReturnDocument(options.Before)
	var doc diskID
	filter := bson.M{"_id": diskIDTblID}
	err := t.tbl.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"disk_id": 1}}, opt).Decode(&doc)
	return doc.DiskID, err
}
