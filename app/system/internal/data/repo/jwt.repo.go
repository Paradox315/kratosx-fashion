package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/cypher"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

const (
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	// jwtBlacklistKey holds the key used to store the JWT Token in the redis.
	jwtBlacklistKey = "jwt:blacklist:%s"

	// jwtBlacklistGracePeriod holds the grace period for the JWT Token in the redis.
	jwtBlacklistGracePeriod = time.Second * 1
)

var (
	ErrMissingJwtToken        = kerrors.Unauthorized(reason, "JWT token is missing")
	ErrTokenInvalid           = kerrors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = kerrors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = kerrors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = kerrors.Unauthorized(reason, "Wrong signing method")
	ErrParseClaimsFail        = kerrors.Unauthorized(reason, "Fail to parse claims")
)

type JwtRepo struct {
	dao *data.Data
	cfg *conf.JWT
	log *log.Helper
}

func NewJwtRepo(dao *data.Data, cfg *conf.JWT, logger log.Logger) biz.JwtRepo {
	return &JwtRepo{
		dao: dao,
		cfg: cfg,
		log: log.NewHelper(log.With(logger, "repo", "jwt")),
	}
}
func (j *JwtRepo) Create(ctx context.Context, user biz.JwtUser) (*biz.Token, error) {
	jti, _ := uuid.NewUUID()
	exp := time.Now().Unix() + j.cfg.Ttl.Seconds
	claims := model.CustomClaims{
		Username: user.GetUsername(),
		Nickname: user.GetNickname(),
		RoleIDs:  user.GetRoleIDs(),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.GetUid(),
			Id:        jti.String(),
			Issuer:    j.cfg.Issuer,
			NotBefore: time.Now().Unix() - 1000,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenStr, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
		err = errors.Wrap(err, "jwt.SignedString")
		j.log.WithContext(ctx).Error("Failed to sign token: %s", err.Error())
		return nil, err
	}
	return &biz.Token{
		AccessToken: tokenStr,
		ExpiresAt:   exp,
		TokenType:   bearerWord,
	}, nil
}

func (j *JwtRepo) IsInBlackList(ctx context.Context, token string) bool {
	joinUnixStr, err := j.dao.RDB.Get(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(token))).Result()
	if err != nil || joinUnixStr == "" {
		return false
	}
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "strconv.ParseInt")
		j.log.WithContext(ctx).Error(err)
		return false
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	return time.Now().Unix()-joinUnix < 1000
}

func (j *JwtRepo) JoinInBlackList(ctx context.Context, token string) error {
	claims, err := j.ParseToken(ctx, token)
	if err != nil {
		return err
	}
	nowUnix := time.Now().Unix()
	timer := time.Duration(claims.ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	return j.dao.RDB.SetEX(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(token)), nowUnix, timer).Err()
}

func (j *JwtRepo) GetSecretKey() string {
	return j.cfg.Secret
}

func (j *JwtRepo) GetIssuer() string {
	return j.cfg.Issuer
}
func (j *JwtRepo) ParseToken(ctx context.Context, token string) (claims *model.CustomClaims, err error) {
	tokenInfo, err := jwt.ParseWithClaims(token, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = ErrTokenInvalid
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				err = ErrTokenExpired
				return
			} else {
				err = ErrTokenParseFail
				return
			}
		}
		err = kerrors.Unauthorized(reason, err.Error())
		return
	} else if !tokenInfo.Valid {
		err = ErrTokenInvalid
		return
	} else if tokenInfo.Method != jwt.SigningMethodHS256 {
		err = ErrUnSupportSigningMethod
		return
	} else if j.IsInBlackList(ctx, token) {
		err = ErrTokenInvalid
		return
	}
	if _, ok := tokenInfo.Claims.(*model.CustomClaims); !ok {
		err = ErrParseClaimsFail
		return
	}
	claims = tokenInfo.Claims.(*model.CustomClaims)
	if claims.Issuer != j.cfg.Issuer {
		err = ErrTokenInvalid
		return
	}
	return
}
