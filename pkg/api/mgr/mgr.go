package mgr

import (
	"github.com/ocean-rw/ocean/pkg/proto"
)

type ListDisksArgs struct {
	Host  string          `in:"query=host"`
	State proto.DiskState `in:"query=state"`
}
