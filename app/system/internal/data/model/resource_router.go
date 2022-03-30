package model

type ResourceRouter struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Ptype  string `json:"ptype" gorm:"column:ptype;comment:策略"`
	RoleID string `json:"role_id" gorm:"column:v0;comment:角色ID"`
	Path   string `json:"path" gorm:"column:v1;comment:路径"`
	Method string `json:"method" gorm:"column:v2;comment:方法"`
}
