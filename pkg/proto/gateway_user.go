package proto

type UserInfo struct {
	UID       uint64     // 用户唯一 ID
	Name      string     // 用户名称
	Password  string     // 用户密码
	AccessKey string     // 用户 AccessKey
	SecretKey string     // 用户 SecretKey
	Status    UserStatus // 用户状态
}

type UserStatus uint8

const (
	// UserStatusNormal 用户状态正常
	UserStatusNormal UserStatus = iota
	// UserStatusDisabled 用户被禁用
	UserStatusDisabled
	// UserStatusMAX 用户状态最大值
	UserStatusMAX
)
