package biz

import (
	"kratosx-fashion/app/system/internal/data/model"
	"mime/multipart"
)

var _ JwtUser = (*User)(nil)

const timeFormat = `2006-01-02 15:04:05`

type Captcha struct {
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captcha_id"`
}

type UserSession struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Location struct {
	Country  string             `json:"country"`
	Region   string             `json:"region"`
	City     string             `json:"city"`
	Position map[string]float32 `json:"position"`
}

type Agent struct {
	Name       string           `json:"name"`
	OS         string           `json:"os"`
	Device     string           `json:"device"`
	DeviceType model.DeviceType `json:"device_type"`
}

type RegisterInfo struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

type User struct {
	Id        string     `json:"id"`
	Username  string     `json:"username"`
	Avatar    string     `json:"avatar"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile"`
	Nickname  string     `json:"nickname"`
	Gender    string     `json:"gender"`
	Status    uint32     `json:"status"`
	Creator   string     `json:"creator"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Roles     []UserRole `json:"roles"`
}

func (u User) GetUid() string {
	return u.Id
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetRoleIDs() (rids []string) {
	for _, role := range u.Roles {
		rids = append(rids, role.Id)
	}
	return
}

func (u User) GetNickname() string {
	return u.Nickname
}

type UserRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UploadInfo struct {
	File *multipart.FileHeader `form:"file"`
}

type SQLOption struct {
	Where string        `json:"query"`
	Order string        `json:"order"`
	Args  []interface{} `json:"args"`
}

type RouterGroup struct {
	Path    string         `json:"path"`
	Name    string         `json:"name"`
	Methods string         `json:"methods"`
	Router  []model.Router `json:"router"`
}

type Menu struct {
	Id        string       `json:"id"`
	ParentId  string       `json:"parent_id"`
	Path      string       `json:"path"`
	Name      string       `json:"name"`
	Component string       `json:"component"`
	Meta      *MenuMeta    `json:"meta"`
	Children  []Menu       `json:"children"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Actions   []MenuAction `json:"actions"`
}

type MenuMeta struct {
	Roles       []string `json:"roles"`
	RequireAuth bool     `json:"require_auth"`
	Icon        string   `json:"icon"`
	Locale      string   `json:"locale"`
	Order       uint32   `json:"order"`
	HideInMenu  bool     `json:"hide_in_menu"`
	NoAffix     bool     `json:"no_affix"`
	IgnoreCache bool     `json:"ignore_cache"`
}

type MenuAction struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
