package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ocean-rw/ocean/internal/master/db/common"
)

var _ common.IDTableIF = (*IDTable)(nil)

const (
	diskIDStart = uint64(10000)
	diskID      = "disk_id"

	stripeIDStart = uint64(10000)
	stripeID      = "stripe_id"
)

type IDTable struct {
	tbl *mongo.Collection
}

func OpenIDTable(tbl *mongo.Collection) (*IDTable, error) {
	if err := initIDTable(tbl, stripeID, stripeIDStart); err != nil {
		return nil, err
	}
	if err := initIDTable(tbl, diskID, diskIDStart); err != nil {
		return nil, err
	}
	return &IDTable{tbl: tbl}, nil
}

func initIDTable(tbl *mongo.Collection, id string, start uint64) error {
	ctx := context.Background()
	cursor, err := tbl.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	var results []bson.M
	err = cursor.All(ctx, &results)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		_, err = tbl.InsertOne(ctx, bson.M{"_id": id, id: start})
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

//func initStripeIDTable(tbl *mongo.Collection) error {
//	ctx := context.Background()
//	cursor, err := tbl.Find(ctx, bson.M{"_id": stripeID})
//	if err != nil {
//		return err
//	}
//	var results []bson.M
//	err = cursor.All(ctx, &results)
//	if err != nil {
//		return err
//	}
//
//	if len(results) == 0 {
//		_, err = tbl.InsertOne(ctx, bson.M{"_id": stripeID, "stripe_id": stripeIDStart})
//		if mongo.IsDuplicateKeyError(err) {
//			err = nil
//		}
//		return err
//	}
//
//	if len(results) != 1 {
//		return fmt.Errorf("unexpected diskIDAllocTable size %d", len(results))
//	}
//
//	return nil
//}
//
//func initDiskIDTable(tbl *mongo.Collection) error {
//	ctx := context.Background()
//	cursor, err := tbl.Find(ctx, bson.M{"_id": diskID})
//	if err != nil {
//		return err
//	}
//	var results []bson.M
//	err = cursor.All(ctx, &results)
//	if err != nil {
//		return err
//	}
//
//	if len(results) == 0 {
//		_, err = tbl.InsertOne(ctx, bson.M{"_id": diskID, "disk_id": diskIDStart})
//		if mongo.IsDuplicateKeyError(err) {
//			err = nil
//		}
//		return err
//	}
//
//	if len(results) != 1 {
//		return fmt.Errorf("unexpected diskIDAllocTable size %d", len(results))
//	}
//
//	return nil
//}

func (t IDTable) AllocStripeID(ctx context.Context, count int) ([]uint64, error) {
	opt := options.FindOneAndUpdate().SetReturnDocument(options.Before)
	var doc = struct {
		StripeID uint64 `bson:"stripe_id"`
	}{}
	filter := bson.M{"_id": stripeID}
	err := t.tbl.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"stripe_id": count}}, opt).Decode(&doc)
	results := make([]uint64, count)
	for i := 0; i < count; i++ {
		results[i] = uint64(i) + doc.StripeID
	}
	return results, err
}

func (t IDTable) AllocDiskID(ctx context.Context) (uint32, error) {
	opt := options.FindOneAndUpdate().SetReturnDocument(options.Before)
	var doc = struct {
		DiskID uint32 `bson:"disk_id"`
	}{}
	filter := bson.M{"_id": diskID}
	err := t.tbl.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"disk_id": 1}}, opt).Decode(&doc)
	return doc.DiskID, err
}
