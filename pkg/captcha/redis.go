package captcha

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"time"
)

const Captcha = "captcha:%s"

func NewDefaultRedisStore(cli *redis.Client, logger log.Logger) *RedisStore {
	return &RedisStore{
		Expiration: time.Second * 180,
		Ctx:        context.Background(),
		rdb:        cli,
		log:        log.NewHelper(log.With(logger, "repo", "captcha")),
	}
}

type RedisStore struct {
	Expiration time.Duration
	Ctx        context.Context
	rdb        *redis.Client
	log        *log.Helper
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.Ctx = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	err := rs.rdb.Set(rs.Ctx, fmt.Sprintf(Captcha, id), value, rs.Expiration).Err()
	if err != nil {
		rs.log.Error("RedisStoreSetError!", zap.Error(err))
		return err
	}
	return nil
}

func (rs *RedisStore) Get(key string, clear bool) string {
	val, err := rs.rdb.Get(rs.Ctx, key).Result()
	if err != nil {
		rs.log.Error("RedisStoreGetError!", zap.Error(err))
		return ""
	}
	if clear {
		err := rs.rdb.Del(rs.Ctx, key).Err()
		if err != nil {
			rs.log.Error("RedisStoreClearError!", zap.Error(err))
			return ""
		}
	}
	return val
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	v := rs.Get(fmt.Sprintf(Captcha, id), clear)
	return v == answer
}
