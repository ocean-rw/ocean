package common

import (
	"context"

	"github.com/ocean-rw/ocean/pkg/api/master"
	"github.com/ocean-rw/ocean/pkg/proto"
)

type Database struct {
	CloseFn func(ctx context.Context) error

	ConfigTable ConfigTableIF
	IDTable     IDTableIF
	DiskTable   DiskTableIF
	StripeTable StripeTableIF
}

type ConfigTableIF interface {
	ClusterID(ctx context.Context) (string, error)
}

type IDTableIF interface {
	AllocStripeID(ctx context.Context, count int) ([]uint64, error)
	AllocDiskID(ctx context.Context) (uint32, error)
}

type DiskTableIF interface {
	Insert(ctx context.Context, disk *proto.Disk) error
	Get(ctx context.Context, diskID uint32) (*proto.Disk, error)
	List(ctx context.Context, args *master.ListDisksArgs) ([]*proto.Disk, error)
	AllocDisks(ctx context.Context, mode proto.Mode, args *master.AllocDisksArgs) ([]uint32, error)
}

type StripeTableIF interface {
	Insert(ctx context.Context, stripes []*proto.Stripe) error
}
