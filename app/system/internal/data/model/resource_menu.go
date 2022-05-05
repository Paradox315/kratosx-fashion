package model

import (
	"github.com/gofiber/fiber/v2/utils"
)

type AuthType uint8

const (
	NotRequireAuth = iota
	RequireAuth
)

type HideType uint8

const (
	ShowInMenu = iota
	HideInMenu
)

type AffixType uint8

const (
	Affix = iota
	NoAffix
)

type CacheType uint8

const (
	Cache = iota
	IgnoreCache
)

type ResourceMenu struct {
	Model
	Name        string    `gorm:"size:50;index;default:'';not null;comment:菜单名称"`         // 菜单名称
	Path        string    `gorm:"size:255;default:'';not null;comment:服务端路由路径URL"`        // 访问路由组(后端)
	Component   string    `gorm:"size:255;default:'';not null;comment:对应前端的路径"`           //对应前端的路径
	Description string    `gorm:"size:255;default:'';not null;comment:菜单描述"`              // 菜单描述
	ParentID    uint      `json:"parent_id" gorm:"index;default:0;not null;comment:父级内码"` // 父级内码
	Locale      string    `gorm:"size:50;default:'';not null;comment:语言"`                 // 语言
	RequireAuth AuthType  `gorm:"default:1;type:tinyint(1);not null;comment:是否需要权限"`      // 是否需要权限
	Icon        string    `gorm:"size:50;default:'';not null;comment:图标"`                 // 图标
	Order       uint32    `gorm:"default:0;type:tinyint(1);not null;comment:排序"`          // 排序
	HideInMenu  HideType  `gorm:"default:0;type:tinyint(1);not null;comment:是否在菜单中隐藏"`    // 是否在菜单中隐藏
	NoAffix     AffixType `gorm:"default:0;type:tinyint(1);not null;comment:是否不附加到菜单"`    // 是否不附加到菜单
	IgnoreCache CacheType `gorm:"default:0;type:tinyint(1);not null;comment:是否忽略缓存"`      // 是否忽略缓存
	Actions     string    `gorm:"size:128;default:'';not null;comment:前端动作按钮ID"`          //动作按钮
}

func (m ResourceMenu) GetActions() (actions []*ResourceAction) {
	_ = codec.Unmarshal(utils.UnsafeBytes(m.Actions), &actions)
	return
}

func (m *ResourceMenu) SetActions(actions []*ResourceAction) {
	bytes, _ := codec.Marshal(actions)
	m.Actions = utils.UnsafeString(bytes)
}
