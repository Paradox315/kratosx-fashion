package biz

import "mime/multipart"

type Captcha struct {
	Captcha   string `json:"captcha,omitempty"`
	CaptchaId string `json:"captcha_id,omitempty"`
}

type UserSession struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
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
