package data

import (
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/logutil"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewRedis,

	NewDiscovery,
	NewRegistrar,
	NewLogger,
	NewStorage,

	NewLoginLogRepo,
	NewUserRepo,
	NewUserRoleRepo,
	NewRoleRepo,
	NewRoleMenuRepo,
	NewMenuRepo,
	NewMenuActionRepo,
	NewMenuActionResourceRepo,
	NewCaptchaRepo,
)

// Data .
type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rdb *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		if err := rdb.Close(); err != nil {
			log.NewHelper(logger).Fatal("redis close error", zap.Error(err))
		}
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		DB:  db,
		RDB: rdb,
	}, cleanup, nil
}

func NewDB(c *conf.Data, logger log.Logger) *gorm.DB {
	zlog := zapgorm2.New(logger.(*logutil.Logger).GetZap())
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger:                                   zlog,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{
			//SingularTable: true, // 表名是否加 s
		},
		PrepareStmt: true,
	})

	if err != nil {
		log.NewHelper(logger).Fatal("failed to connect database", zap.Error(err))
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	sqlDb.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	if c.Database.AutoMigrate {
		if err = db.AutoMigrate(
			&model.LoginLog{},
			&model.User{},
			&model.UserRole{},
			&model.Role{},
			&model.RoleMenu{},
			&model.Menu{},
			&model.MenuAction{},
			&model.MenuActionResource{},
		); err != nil {
			panic(err)
		}
	}
	return db
}

func NewRedis(c *conf.Data, logger log.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	return rdb
}
