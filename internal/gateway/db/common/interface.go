package common

import (
	"context"

	"github.com/ocean-rw/ocean/pkg/proto"
)

type Database struct {
	CloseFn func(ctx context.Context) error

	UserTable   UserTableIF
	BucketTable BucketTableIF
	FileTable   FileTableIF
}

type UserTableIF interface {
}

type BucketTableIF interface {
	List(ctx context.Context, uid uint64) ([]*proto.BucketInfo, error)
	Insert(ctx context.Context, bucketInfo *proto.BucketInfo) error
	Get(ctx context.Context, bid string) (*proto.BucketInfo, error)
	Delete(ctx context.Context, bid string) error
}

type FileTableIF interface {
	Upsert(ctx context.Context, fileInfo *proto.FileInfo) error
	Get(ctx context.Context, bid string, fid string) (*proto.FileInfo, error)
	List(ctx context.Context, bid string) ([]*proto.FileInfo, error)
}
