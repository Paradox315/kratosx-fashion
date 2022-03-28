package model

import (
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/middleware"
	"strconv"
)

var _ middleware.JwtUser = (*User)(nil)

type UserStatus uint8
type GenderStatus uint8

const (
	UserStatusNormal UserStatus = iota + 1
	UserStatusForbid
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusNormal:
		return "正常"
	case UserStatusForbid:
		return "禁用"
	default:
		return "未知"
	}
}

const (
	GenderUnknown GenderStatus = iota
	GenderFemale
	GenderMale
)

func (g GenderStatus) String() string {

	switch g {
	case GenderUnknown:
		return "未知"
	case GenderFemale:
		return "女"
	case GenderMale:
		return "男"
	default:
		return "未知"
	}
}

// User
// 用户表
type User struct {
	gorm.Model
	Username  string       `gorm:"size:64;uniqueIndex;default:'unknown';not null;comment:用户名"`                                                     // 用户名
	Nickname  string       `gorm:"size:64;default:'';not null;comment:昵称"`                                                                         // 昵称
	Password  string       `gorm:"size:255;default:'';not null;comment:密码"`                                                                        // 密码
	Avatar    string       `gorm:"size:255;default:'https://paradox-hyw.oss-cn-shanghai.aliyuncs.com/img/default-avatar.png';not null;comment:头像"` // 头像
	Email     string       `gorm:"size:255;uniqueIndex;default:'';not null;comment:邮箱"`                                                            // 邮箱
	Mobile    string       `gorm:"size:20;uniqueIndex;default:'';not null;comment:手机号"`                                                            // 手机号
	Gender    GenderStatus `gorm:"default:0;type:tinyint(1);not null;comment:状态(0:未知 1:男性 2:女性)"`
	Status    UserStatus   `gorm:"default:1;type:tinyint(1);not null;comment:状态(1:启用 2:停用)"` // 状态(1:启用 2:停用)
	Creator   string       `gorm:"default:'admin';not null;comment:创建者"`                     // 创建者
	Extras    string       `gorm:"size:1024;default:'';not null;comment:扩展字段信息"`             //其他字段
	Roles     []Role       `gorm:"many2many:user_roles;"`
	LoginLogs []LoginLog   `gorm:"foreignkey:UserID"`
}

func (u User) GetUid() string {
	return strconv.Itoa(int(u.ID))
}
