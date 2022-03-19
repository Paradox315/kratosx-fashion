package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName  string     `gorm:"size:64;uniqueIndex;default:'unknown';not null;comment:用户名"`     // 用户名
	NickName  string     `gorm:"size:64;index;default:'';not null;comment:昵称"`                   // 昵称
	Password  string     `gorm:"size:40;default:'';not null;comment:密码"`                         // 密码
	Email     string     `gorm:"size:255;default:'';not null;comment:邮箱"`                        // 邮箱
	Phone     string     `gorm:"size:20;default:'';not null;comment:手机号"`                        // 手机号
	Status    uint8      `gorm:"index;default:0;type:tinyint(1);not null;comment:状态(1:启用 2:停用)"` // 状态(1:启用 2:停用)
	CreatorID uint64     `gorm:"default:0;not null;comment:创建者ID"`                               // 创建者
	Extras    string     `gorm:"size:1024;default:'';not null;comment:扩展字段信息"`                   //其他字段
	Roles     []Role     `gorm:"many2many:user_roles;"`
	LoginLogs []LoginLog `gorm:"foreignkey:UserID"`
}
