package model

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"time"
)

var (
	codec = encoding.GetCodec("json")
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}
