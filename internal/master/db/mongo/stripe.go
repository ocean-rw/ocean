package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/master/db/common"
	"github.com/ocean-rw/ocean/pkg/proto"
)

var _ common.StripeTableIF = (*StripeTable)(nil)

type StripeTable struct {
	tbl *mongo.Collection
}

func OpenStripeTable(tbl *mongo.Collection) (*StripeTable, error) {
	return &StripeTable{tbl: tbl}, nil
}

func (t *StripeTable) Insert(ctx context.Context, stripes []*proto.Stripe) error {
	data := make([]interface{}, len(stripes))
	for i := range stripes {
		data[i] = stripes[i]
	}
	_, err := t.tbl.InsertMany(ctx, data)
	return err
}
