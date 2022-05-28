package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
)

var _ common.UserTableIF = (*userTable)(nil)

type userTable struct {
	tbl *mongo.Collection
}

func openUserTable(tbl *mongo.Collection) (*userTable, error) {
	return &userTable{tbl: tbl}, nil
}
