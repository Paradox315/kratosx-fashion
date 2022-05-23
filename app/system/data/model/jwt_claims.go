package model

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	RoleIDs  []uint `json:"role_ids"`
	UID      uint   `json:"uid"`
	jwt.StandardClaims
}
