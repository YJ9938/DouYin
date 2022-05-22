package model

// 用户信息表
type User struct {
	ID       int64  `gorm:"primary key" json:"id"`
	Username string `gorm:"not null; unique; size:32; index:idx_username" json:"username"`
	Password []byte `gorm:"not null; type:varbinary(256)" json:"password"`
	Salt     []byte `gorm:"not null; type:varbinary(32)" json:"salt"`
}
