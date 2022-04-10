package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/pkg/xsync"
	"os"
	"strings"
	"sync"
	"time"

	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
)

var _ kmw.FiberMiddleware = (*JWTService)(nil)

const (
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	// jwtBlacklistKey holds the key used to store the JWT Token in the redis.
	jwtBlacklistKey = "jwt:blacklist:%s"

	// jwtBlacklistGracePeriod holds the grace period for the JWT Token in the redis.
	jwtBlacklistGracePeriod = time.Second * 1
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(reason, "JWT token is missing")
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrParseClaimsFail        = errors.Unauthorized(reason, "Fail to parse claims")
)

type JWTService struct {
	jwtRepo biz.JwtRepo
	uc      *biz.UserUsecase
	once    sync.Once
	rdb     *redis.Client
	log     *log.Helper
	lock    xsync.XMutex
}

func NewJwtService(jwtRepo biz.JwtRepo, uc *biz.UserUsecase, rdb *redis.Client, logger log.Logger) *JWTService {
	j := &JWTService{
		jwtRepo: jwtRepo,
		uc:      uc,
		rdb:     rdb,
		log:     log.NewHelper(logger),
		lock:    xsync.Lock("refresh_token_lock", 2000, rdb),
	}
	j.once.Do(func() {
		kmw.RegisterMiddleware(j)
	})
	return j
}

func (j *JWTService) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if os.Getenv("JWT") == "false" {
			return c.Next()
		}
		errCatch := func(ctx context.Context) error {
			auths := strings.SplitN(c.Get(authorizationKey), " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
				return ErrMissingJwtToken
			}
			jwtToken := auths[1]
			claims, err := j.jwtRepo.ParseToken(ctx, jwtToken)
			if err != nil {
				return err
			}
			if claims.ExpiresAt-time.Now().Unix() < jwtBlacklistGracePeriod.Milliseconds() {
				if j.lock.Get() {
					var user biz.JwtUser
					user, err = j.uc.Get(ctx, cast.ToUint(claims.Id))
					if err != nil {
						j.log.WithContext(ctx).Error("get user info error", err)
						j.lock.Release()
					} else {
						tokenData, _ := j.jwtRepo.Create(ctx, user)
						c.Set("new-token", tokenData.AccessToken)
						c.Set("new-expires-at", cast.ToString(tokenData.ExpiresAt))
						_ = j.jwtRepo.JoinInBlackList(ctx, jwtToken)
					}
				}
			}
			c.Locals("uid", claims.UID)
			c.Locals("username", claims.Username)
			c.Locals("roles", claims.RoleIDs)
			c.Locals("nickname", claims.Nickname)
			c.Locals("token", jwtToken)
			return nil
		}
		if err := errCatch(c.Context()); err != nil {
			j.log.WithContext(c.Context()).Error(err)
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return c.Next()
	}
}

func (j *JWTService) Name() string {
	return kmw.AuthenticatorCfg
}
