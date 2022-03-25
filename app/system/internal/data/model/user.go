package model

import (
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/middleware"
	"strconv"
)

var _ middleware.JwtUser = (*User)(nil)

type User struct {
	gorm.Model
	Username  string     `gorm:"size:64;uniqueIndex;default:'unknown';not null;comment:用户名"`                                                     // 用户名
	Nickname  string     `gorm:"size:64;default:'';not null;comment:昵称"`                                                                         // 昵称
	Password  string     `gorm:"size:255;default:'';not null;comment:密码"`                                                                        // 密码
	Avatar    string     `gorm:"size:255;default:'https://paradox-hyw.oss-cn-shanghai.aliyuncs.com/img/default-avatar.png';not null;comment:头像"` // 头像
	Email     string     `gorm:"size:255;uniqueIndex;default:'';not null;comment:邮箱"`                                                            // 邮箱
	Mobile    string     `gorm:"size:20;uniqueIndex;default:'';not null;comment:手机号"`                                                            // 手机号
	Status    uint8      `gorm:"index;default:1;type:tinyint(1);not null;comment:状态(1:启用 2:停用)"`                                                 // 状态(1:启用 2:停用)
	CreatorID uint64     `gorm:"default:0;not null;comment:创建者ID"`                                                                               // 创建者
	Extras    string     `gorm:"size:1024;default:'';not null;comment:扩展字段信息"`                                                                   //其他字段
	Roles     []Role     `gorm:"many2many:user_roles;"`
	LoginLogs []LoginLog `gorm:"foreignkey:UserID"`
}

func (u User) GetUid() string {
	return strconv.Itoa(int(u.ID))
}

func (u User) GetCreatorID() string {
	return strconv.Itoa(int(u.CreatorID))
}
