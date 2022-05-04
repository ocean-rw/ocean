package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/ocean-mgr/db/common"
)

var _ common.StripeTableIF = (*StripeTable)(nil)

type StripeTable struct {
	tbl *mongo.Collection
}

func OpenStripeTable(tbl *mongo.Collection) (*StripeTable, error) {
	return &StripeTable{tbl: tbl}, nil
}
