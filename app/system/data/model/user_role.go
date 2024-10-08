package model

// UserRole
// 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primarykey"`
	UserID uint `gorm:"index;default:0;not null;comment:用户ID"` // 用户内码
	RoleID uint `gorm:"index;default:0;not null;comment:角色ID"` // 角色内码
}
