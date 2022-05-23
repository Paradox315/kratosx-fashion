package model

type ResourceRouter struct {
	Ptype  string `json:"ptype" gorm:"column:ptype;comment:策略"`
	RoleID string `json:"role_id" gorm:"column:v0;comment:角色ID"`
	Path   string `json:"path" gorm:"column:v1;comment:路径"`
	Method string `json:"method" gorm:"column:v2;comment:方法"`
}
type Router struct {
	Method string   `json:"method"`
	Path   string   `json:"path"`
	Name   string   `json:"name"`
	Params []string `json:"params"`
	Group  string   `json:"group"`
}
