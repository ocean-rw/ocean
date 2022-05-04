package common

import (
	"context"

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
	AllocDiskID(ctx context.Context) (uint32, error)
}

type DiskTableIF interface {
	Insert(ctx context.Context, disk *proto.Disk) error
	Get(ctx context.Context, diskID uint32) (*proto.Disk, error)
	List(ctx context.Context, args *proto.ListDisksArgs) ([]*proto.Disk, error)
}

type StripeTableIF interface {
}
