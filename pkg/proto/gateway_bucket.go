package proto

import (
	"time"
)

type BucketInfo struct {
	BID     string       `bson:"_id"`      // 桶 ID
	UID     uint64       `bson:"uid"`      // 归属用户 ID
	PutTime time.Time    `bson:"put_time"` // 创建时间
	DelTime time.Time    `bson:"del_time"` // 删除时间
	Status  BucketStatus `bson:"status"`   // 桶状态
}

type BucketStatus uint8

const (
	// BucketStatusNormal 桶状态正常
	BucketStatusNormal BucketStatus = iota
	// BucketStatusDisabled 桶被禁用
	BucketStatusDisabled
	// BucketStatusMarkDeleted 桶被标记删除
	BucketStatusMarkDeleted
	// BucketStatusReclaimed 桶被删除且数据被删除
	BucketStatusReclaimed
	// BucketStatusMAX 桶状态最大值
	BucketStatusMAX
)
