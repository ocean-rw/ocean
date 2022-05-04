package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func findAll(ctx context.Context, tbl *mongo.Collection, results interface{}) error {
	cursor, err := tbl.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	err = cursor.All(ctx, results)
	if err != nil {
		return err
	}
	return nil
}
