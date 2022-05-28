package proto

import (
	"time"
)

type FileInfo struct {
	FID      string            `bson:"_id"`       // 名称
	BID      string            `bson:"bid"`       // 归属 bucket
	Size     int64             `bson:"size"`      // 文件大小
	Status   FileStatus        `bson:"status"`    // 状态
	PutTime  time.Time         `bson:"put_time"`  // 上传时间
	DelTime  time.Time         `bson:"del_time"`  // 标删时间
	Hash     string            `bson:"hash"`      // 文件 Hash，MD5 或者 SHA1
	ETag     string            `bson:"etag"`      // 文件 ETag
	MimeType string            `bson:"mime_type"` // 媒体类型
	Meta     map[string]string `bson:"meta"`      // 用户自定义元数据，kv 形式
	FD       string            `bson:"fd"`        // 文件位置索引
}

type FileStatus uint8

const (
	// FileStatusNormal 文件状态正常
	FileStatusNormal FileStatus = iota
	// FileStatusDisabled 文件被禁用
	FileStatusDisabled
	// FileStatusMarkDeleted 表示文件被标记删除
	FileStatusMarkDeleted
	// FileStatusReclaimed 文件被回收
	FileStatusReclaimed
	// FileStatusMAX 文件状态最大值
	FileStatusMAX
)
