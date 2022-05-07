package proto

import (
	"errors"
	"time"
)

// ErrUnknownDiskMode unknown disk mode
var ErrUnknownDiskMode = errors.New("unknown disk mode")

// DiskState 表示磁盘的状态。
type DiskState uint8

const (
	// DiskStateNormal 表示磁盘正常
	DiskStateNormal DiskState = iota + 1
	// DiskStateReadonly 表示磁盘只读
	DiskStateReadonly
	// DiskStateBroken 表示磁盘已经损坏
	DiskStateBroken
	// DiskStateRepairing 表示磁盘中的数据在修复中
	DiskStateRepairing
	// DiskStateRepaired 表示磁盘中的数据已经修复完成
	DiskStateRepaired
	// DiskStateMAX 表示最大值
	DiskStateMAX
)

type Disk struct {
	*DiskLabel

	Host      string    `json:"host" bson:"host"`
	Path      string    `json:"path" bson:"path"`
	Capacity  uint64    `json:"capacity" bson:"capacity"`
	Available uint64    `json:"available" bson:"available"`
	State     DiskState `json:"state" bson:"state"`
	CreateAt  time.Time `json:"create_at" bson:"create_at"`
	ModifyAt  time.Time `json:"modify_at" bson:"modify_at"`
}

type DiskLabel struct {
	ClusterID string `json:"cluster_id" bson:"cluster_id"`
	DiskID    uint32 `json:"disk_id" bson:"disk_id"`
}
