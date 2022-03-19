package infra

import (
	"github.com/jassue/go-storage/local"
	"github.com/jassue/go-storage/oss"
	"github.com/jassue/go-storage/storage"
	"github.com/jinzhu/copier"
	"kratosx-fashion/app/system/internal/conf"
)

func NewStorage(sc *conf.Storage) storage.Storage {
	switch sc.Type {
	case string(storage.Local):
		var localCfg local.Config
		if err := copier.Copy(&localCfg, sc.Disks.Local); err != nil {
			panic(err)
		}
		disk, err := local.Init(localCfg)
		if err != nil {
			panic(err)
		}
		return disk
	case string(storage.Oss):
		var ossCfg oss.Config
		if err := copier.Copy(&ossCfg, sc.Disks.AliOss); err != nil {
			panic(err)
		}
		disk, err := oss.Init(ossCfg)
		if err != nil {
			panic(err)
		}
		return disk
	default:
		panic("storage type not support")
	}
}
