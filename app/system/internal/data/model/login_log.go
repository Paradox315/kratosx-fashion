package model

import "gorm.io/gorm"

type LoginType uint8

const (
	LoginType_Login LoginType = iota + 1
	LoginType_Logout
)

func (l LoginType) String() string {
	switch l {
	case LoginType_Login:
		return "登录"
	case LoginType_Logout:
		return "退出"
	default:
		return "未知"
	}
}

type DeviceType uint8

const (
	DeviceType_PC DeviceType = iota + 1
	DeviceType_Mobile
	DeviceType_Pad
	DeviceType_Bot
)

func (d DeviceType) String() string {
	switch d {
	case DeviceType_PC:
		return "PC"
	case DeviceType_Mobile:
		return "移动端"
	case DeviceType_Pad:
		return "平板"
	case DeviceType_Bot:
		return "机器人"
	default:
		return "未知"
	}
}

type LoginLog struct {
	gorm.Model
	UserID     uint64 `gorm:"index;type:unsigned bigint;default:0;not null;comment:'用户ID'"`
	Ip         string `gorm:"type:varchar(40);default:'';not null;comment:'IP地址'"`
	Location   string `gorm:"type:varchar(40);default:'';not null;comment:'地理位置'"`
	LoginType  uint8  `gorm:"type:tinyint(1);default:1;not null;comment:'登录类型：1-登录，2-退出'"`
	Agent      string `gorm:"type:varchar(255);default:'';not null;comment:'浏览器'"`
	OS         string `gorm:"type:varchar(40);default:'';not null;comment:'操作系统'"`
	Device     string `gorm:"type:varchar(40);default:'';not null;comment:'设备'"`
	DeviceType uint8  `gorm:"type:tinyint(1);default:0;not null;comment:'设备类型：1-PC，2-移动，3-平板，4-爬虫'"`
}
