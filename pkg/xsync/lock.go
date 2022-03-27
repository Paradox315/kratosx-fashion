package xsync

import (
	"context"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

type XMutex interface {
	Get() bool
	Block(seconds int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string // 锁名称
	owner   string // 锁标识
	seconds int64  // 有效期
	rdb     *redis.Client
}

// 释放锁 Lua 脚本，防止任何客户端都能解锁
const releaseLockLuaScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`

// Lock 生成锁
func Lock(name string, seconds int64, rdb *redis.Client) XMutex {
	return &lock{
		context: context.Background(),
		name:    name,
		owner:   RandString(16),
		seconds: seconds,
		rdb:     rdb,
	}
}

// Get 获取锁
func (l *lock) Get() bool {
	return l.rdb.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// Block 阻塞一段时间，尝试获取锁
func (l *lock) Block(seconds int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-seconds >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// Release 释放锁
func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaScript)
	result := luaScript.Run(l.context, l.rdb, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

// ForceRelease 强制释放锁
func (l *lock) ForceRelease() {
	l.rdb.Del(l.context, l.name).Val()
}
