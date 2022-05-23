package model

import "github.com/gofiber/fiber/v2/utils"

type DeviceType uint8

const (
	DevicetypePc DeviceType = iota + 1
	DevicetypeMobile
	DevicetypePad
	DevicetypeBot
)

func (d DeviceType) String() string {
	switch d {
	case DevicetypePc:
		return "PC"
	case DevicetypeMobile:
		return "移动端"
	case DevicetypePad:
		return "平板"
	case DevicetypeBot:
		return "机器人"
	default:
		return "未知"
	}
}

type UserLog struct {
	Model
	UID        uint       `gorm:"index;default:0;not null;comment:'用户ID'"`
	Ip         string     `gorm:"type:varchar(40);default:'';not null;comment:'IP地址'"`
	Method     string     `gorm:"type:varchar(20);default:'';not null;comment:'请求方式'"`
	Path       string     `gorm:"type:varchar(128);default:'';not null;comment:'请求路径'"`
	Status     uint32     `gorm:"default:0;not null;comment:'状态码'"`
	Country    string     `gorm:"type:varchar(30);default:'';not null;comment:'国家'"`
	Region     string     `gorm:"type:varchar(30);default:'';not null;comment:'地区'"`
	City       string     `gorm:"type:varchar(30);default:'';not null;comment:'城市'"`
	Position   string     `gorm:"type:varchar(128);default:'';not null;comment:'位置'"`
	UserAgent  string     `gorm:"type:varchar(255);default:'';not null;comment:'浏览器'"`
	Client     string     `gorm:"type:varchar(60);default:'';not null;comment:'客户端'"`
	OS         string     `gorm:"type:varchar(40);default:'';not null;comment:'操作系统'"`
	Device     string     `gorm:"type:varchar(40);default:'';not null;comment:'设备'"`
	DeviceType DeviceType `gorm:"type:tinyint(1);default:0;not null;comment:'设备类型：1-PC，2-移动，3-平板，4-爬虫'"`
	Type       string     `gorm:"type:varchar(20);default:'';not null;comment:'类型'"`
}

func (u UserLog) GetPosition() (pos map[string]float32) {
	_ = codec.Unmarshal(utils.UnsafeBytes(u.Position), &pos)
	return
}

func (u *UserLog) SetPosition(pos map[string]float32) {
	bytes, _ := codec.Marshal(pos)
	u.Position = utils.UnsafeString(bytes)
}
