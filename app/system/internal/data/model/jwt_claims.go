package model

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	RoleIDs  []string `json:"role_ids"`
	UID      string   `json:"uid"`
	jwt.StandardClaims
}
