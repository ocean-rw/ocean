package common

import (
	"context"
	"io"

	"github.com/ocean-rw/ocean/pkg/proto"
)

type DiskIF interface {
	Init(context.Context, *proto.DiskLabel) error
	Stat(context.Context) (*proto.Disk, error)
	Put(context.Context, string, io.ReadCloser) error
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
}
