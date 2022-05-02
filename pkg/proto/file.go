package proto

import (
	"time"
)

type FileInfo struct {
	BID      uint64            // 归属 bucket id
	Name     string            // 名称
	Status   FileStatus        // 状态
	PutTime  time.Time         // 上传时间
	DelTime  time.Time         // 标删时间
	Hash     string            // 文件 Hash，MD5 或者 SHA1
	MimeType string            // 媒体类型
	Meta     map[string]string // 用户自定义元数据，kv 形式
	FD       string            // 文件位置索引
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
