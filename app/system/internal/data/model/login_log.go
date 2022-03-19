package model

import "gorm.io/gorm"

type LoginLog struct {
	gorm.Model
	UserID    uint64 `gorm:"column:user_id;type:unsigned bigint;default:0; not null"`
	Ip        string `gorm:"column:ip;type:varchar(40);not null"`
	Location  string `gorm:"column:location;type:varchar(40);not null"`
	LoginType uint8  `gorm:"column:login_type;type:tinyint(1);default:0;not null"`
	Agent     string `gorm:"column:agent;type:varchar(255);not null"`
	OS        string `gorm:"column:os;type:varchar(40);not null"`
}
