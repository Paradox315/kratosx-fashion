package middleware

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/pkg/cypher"
	"kratosx-fashion/pkg/xsync"
	"strconv"
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

type JwtUser interface {
	GetUid() string
}

type TokenOut struct {
	Token     string `json:"token"`
	ExpireAt  int64  `json:"expire_at"`
	TokenType string `json:"token_type"`
}

type CustomClaims struct {
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	RoleIDs  []string `json:"role_ids"`
	jwt.StandardClaims
}

type JWTService struct {
	secret string
	ttl    int64
	issuer string

	once *sync.Once
	rdb  *redis.Client
	log  *log.Helper
	lock xsync.XMutex
}

func NewJwtService(jc *conf.JWT, rdb *redis.Client, logger log.Logger) *JWTService {
	j := &JWTService{
		secret: jc.Secret,
		ttl:    jc.Ttl.Seconds,
		issuer: jc.Issuer,
		rdb:    rdb,
		once:   new(sync.Once),
		log:    log.NewHelper(logger),
		lock:   xsync.Lock("refresh_token_lock", 2000, rdb),
	}
	j.once.Do(func() {
		kmw.RegisterMiddleware(j)
	})
	return j
}

func (j *JWTService) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		errCatch := func(ctx context.Context) error {
			auths := strings.SplitN(c.Get(authorizationKey), " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
				return ErrMissingJwtToken
			}
			jwtToken := auths[1]
			tokenInfo, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(j.secret), nil
			})
			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						return ErrTokenInvalid
					} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						return ErrTokenExpired
					} else {
						return ErrTokenParseFail
					}
				}
				return errors.Unauthorized(reason, err.Error())
			} else if !tokenInfo.Valid {
				return ErrTokenInvalid
			} else if tokenInfo.Method != jwt.SigningMethodHS256 {
				return ErrUnSupportSigningMethod
			} else if j.InBlackList(c.Context(), jwtToken) {
				return ErrTokenInvalid
			}

			if claims, ok := tokenInfo.Claims.(jwt.StandardClaims); !ok {
				return ErrParseClaimsFail
			} else if claims.Issuer != j.issuer {
				return ErrTokenInvalid
			} else {
				if claims.ExpiresAt-time.Now().Unix() < jwtBlacklistGracePeriod.Milliseconds() {
					if j.lock.Get() {
						//err, user := services.JwtService.GetUserInfo(GuardName, claims.Id)
						//if err != nil {
						//	j.log.WithContext(ctx).Error("get user info error", err)
						//	j.lock.Release()
						//} else {
						//	tokenData, _ := j.CreateToken(ctx, user)
						//	c.Header("new-token", tokenData.AccessToken)
						//	c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
						//	_ = j.JoinBlackList(jwtToken)
						//}
					}
				}
			}

			c.Locals("token", jwtToken)
			return nil
		}
		if err := errCatch(c.Context()); err != nil {
			j.log.WithContext(c.Context()).Error(err)
			return apistate.Error().WithError(err).Send(c)
		}
		return c.Next()
	}
}

func (j *JWTService) Name() string {
	return kmw.AuthenticatorCfg
}

// CreateToken 生成 Token
func (j *JWTService) CreateToken(ctx context.Context, user JwtUser) (tokenData TokenOut, err error) {
	jti, _ := uuid.NewUUID()
	exp := time.Now().Unix() + j.ttl
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: exp,
			Subject:   user.GetUid(),
			Issuer:    j.issuer, // 用于在中间件中区分不同客户端颁发的 token，避免 token 跨端使用
			NotBefore: time.Now().Unix() - 1000,
			IssuedAt:  time.Now().Unix(),
			Id:        jti.String(),
		},
	)

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		j.log.WithContext(ctx).Error("Failed to sign token: %s", err.Error())
		return
	}
	tokenData = TokenOut{
		Token:     tokenStr,
		ExpireAt:  exp,
		TokenType: bearerWord,
	}
	return
}

// JoinBlackList token 加入黑名单
func (j *JWTService) JoinBlackList(ctx context.Context, token string) (err error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	nowUnix := time.Now().Unix()
	timer := time.Duration(parsedToken.Claims.(*jwt.StandardClaims).ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	return j.rdb.SetEX(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(parsedToken.Raw)), nowUnix, timer).Err()
}

func (j *JWTService) InBlackList(ctx context.Context, token string) bool {
	joinUnixStr, err := j.rdb.Get(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(token))).Result()
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		j.log.WithContext(ctx).Error(err)
		return false
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	return time.Now().Unix()-joinUnix < 10
}
