package biz

import (
	"kratosx-fashion/app/system/internal/data/model"
	"mime/multipart"
)

const timeFormat = `2006-01-02 15:04:05`

type Captcha struct {
	Captcha   string `json:"captcha,omitempty"`
	CaptchaId string `json:"captcha_id,omitempty"`
}

type UserSession struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Location struct {
	Country  string             `json:"country,omitempty"`
	Region   string             `json:"region,omitempty"`
	City     string             `json:"city,omitempty"`
	Position map[string]float32 `json:"position,omitempty"`
}

type Agent struct {
	Name       string           `json:"name,omitempty"`
	OS         string           `json:"os,omitempty"`
	Device     string           `json:"device,omitempty"`
	DeviceType model.DeviceType `json:"device_type,omitempty"`
}

type RegisterInfo struct {
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpireAt    int64  `json:"expire_at,omitempty"`
}

type UploadInfo struct {
	File *multipart.FileHeader `form:"file"`
}

type SQLOption struct {
	Where string        `json:"query,omitempty"`
	Order string        `json:"order,omitempty"`
	Args  []interface{} `json:"args,omitempty"`
}
