package mongo

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocean-rw/ocean/internal/master/db/common"
	"github.com/ocean-rw/ocean/pkg/api/master"
	"github.com/ocean-rw/ocean/pkg/proto"
)

var _ common.DiskTableIF = (*DiskTable)(nil)

type DiskTable struct {
	tbl *mongo.Collection
}

func OpenDiskTable(tbl *mongo.Collection) (*DiskTable, error) {
	return &DiskTable{tbl: tbl}, nil
}

func (t *DiskTable) Insert(ctx context.Context, disk *proto.Disk) error {
	disk.CreateAt = time.Now()
	_, err := t.tbl.InsertOne(ctx, disk)
	return err
}

func (t *DiskTable) Get(ctx context.Context, diskID uint32) (*proto.Disk, error) {
	disk := new(proto.Disk)
	err := t.tbl.FindOne(ctx, bson.M{"_id": diskID}).Decode(&disk)
	if err != nil {
		return nil, err
	}
	return disk, err
}

func (t *DiskTable) List(ctx context.Context, args *master.ListDisksArgs) ([]*proto.Disk, error) {
	filter := bson.M{}
	if args != nil {
		if args.Host != "" {
			filter["host"] = args.Host
		}
		if args.State != 0 {
			filter["state"] = args.State
		}
	}

	cursor, err := t.tbl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	disks := make([]*proto.Disk, 0)
	err = cursor.All(ctx, &disks)
	if err != nil {
		return nil, err
	}
	return disks, nil
}

func (t *DiskTable) AllocDisks(ctx context.Context, mode proto.Mode, args *master.AllocDisksArgs) ([]uint32, error) {
	disks, err := t.List(ctx, &master.ListDisksArgs{State: proto.DiskStateNormal})
	if err != nil {
		return nil, err
	}
	num := mode.N() + mode.M()
	total := len(disks)
	if total < num {
		return nil, errors.New("not enough disk")
	}

	exist := make([]bool, total)
	ids := make([]uint32, 0, num)
	for len(ids) < num {
		i := rand.Intn(total)
		if exist[i] {
			continue
		}
		ids = append(ids, disks[i].DiskID)
		exist[i] = true
	}

	return ids, nil
}
