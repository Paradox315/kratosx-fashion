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

	jwtBlacklistKey = "jwt:blacklist:%s"
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

type JWTService struct {
	secret string
	ttl    int64
	issuer string

	once *sync.Once
	rdb  *redis.Client
	log  *log.Helper
}

func NewJwtService(jc *conf.JWT, rdb *redis.Client, logger log.Logger) kmw.FiberMiddleware {
	j := &JWTService{
		secret: jc.Secret,
		ttl:    jc.Ttl.Seconds,
		issuer: jc.Issuer,
		rdb:    rdb,
		once:   new(sync.Once),
		log:    log.NewHelper(logger),
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
			}
			c.Set("token", jwtToken)
			c.Locals("claims", tokenInfo.Claims)
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
func (j *JWTService) CreateToken(ctx context.Context, user JwtUser) (tokenStr string, err error, token *jwt.Token) {
	jti, _ := uuid.NewUUID()
	exp := time.Now().Unix() + j.ttl
	token = jwt.NewWithClaims(
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

	tokenStr, err = token.SignedString([]byte(j.secret))
	if err != nil {
		j.log.WithContext(ctx).Error("Failed to sign token: %s", err.Error())
		return
	}

	return
}

// JoinBlackList token 加入黑名单
func (j *JWTService) JoinBlackList(ctx context.Context, token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*jwt.StandardClaims).ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	return j.rdb.SetEX(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(token.Raw)), nowUnix, timer).Err()
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
