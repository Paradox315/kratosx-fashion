package model

type ResourceAction struct {
	Name string `gorm:"size:100;default:'';not null;comment:动作名称"` // 动作名称
	Code string `gorm:"size:100;default:'';not null;comment:动作编码"` // 动作编码
}
