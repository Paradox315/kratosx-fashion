package model

import "gorm.io/gorm"

type HiddenStatus uint8

const (
	HiddenStatusHide HiddenStatus = iota
	HiddenStatusShow
)

func (h HiddenStatus) String() string {
	switch h {
	case HiddenStatusShow:
		return "show"
	case HiddenStatusHide:
		return "hide"
	default:
		return "unknown"
	}
}

type KeepAliveStatus uint8

const (
	KeepAliveStatusClose KeepAliveStatus = iota
	KeepAliveStatusOpen
)

func (k KeepAliveStatus) String() string {
	switch k {
	case KeepAliveStatusOpen:
		return "open"
	case KeepAliveStatusClose:
		return "close"
	default:
		return "unknown"
	}
}

type RequireAuthStatus uint8

const (
	RequireAuthStatusClose RequireAuthStatus = iota
	RequireAuthStatusOpen
)

func (r RequireAuthStatus) String() string {
	switch r {
	case RequireAuthStatusOpen:
		return "open"
	case RequireAuthStatusClose:
		return "close"
	default:
		return "unknown"
	}
}

type ResourceMenu struct {
	gorm.Model
	Name      string `gorm:"size:50;index;default:'';not null;comment:菜单名称"`         // 菜单名称
	Path      string `gorm:"size:255;default:'';not null;comment:服务端路由路径URL"`        // 访问路由组(后端)
	Component string `gorm:"size:255;default:'';not null;comment:对应前端的路径"`           //对应前端的路径
	ParentID  uint   `json:"parent_id" gorm:"index;default:0;not null;comment:父级内码"` // 父级内码
	MenuMeta  `gorm:"comment:菜单附加属性"`
	Hidden    HiddenStatus    `gorm:"default:0;type:tinyint(1);not null;comment:是否显示"`       // 是否显示
	Keepalive KeepAliveStatus `gorm:"default:0;type:tinyint(1);not null;comment:是否缓存"`       //是否缓存
	Actions   string          `gorm:"size:1024;default:'';not null;comment:前端动作按钮列表，json存储"` //动作按钮
}

type MenuMeta struct {
	Locale      string            `gorm:"size:50;default:'';not null;comment:语言"`            // 语言
	RequireAuth RequireAuthStatus `gorm:"default:1;type:tinyint(1);not null;comment:是否需要权限"` // 是否需要权限
	Icon        string            `gorm:"size:50;default:'';not null;comment:图标"`            // 图标
	Order       uint32            `gorm:"default:0;not null;comment:排序"`                     // 排序
}
