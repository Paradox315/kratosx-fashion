package repo

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/cypher"
	"strconv"
	"time"
)

const (
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// jwtBlacklistKey holds the key used to store the JWT Token in the redis.
	jwtBlacklistKey = "jwt:blacklist:%s"

	// jwtBlacklistGracePeriod holds the grace period for the JWT Token in the redis.
	jwtBlacklistGracePeriod = time.Second * 1
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
	exp := time.Now().Unix() + int64(j.cfg.Ttl.Nanos/1e6)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		model.CustomClaims{
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
		},
	)

	tokenStr, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
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
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		j.log.WithContext(ctx).Error(err)
		return false
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	return time.Now().Unix()-joinUnix < 1000
}

func (j *JwtRepo) JoinInBlackList(ctx context.Context, token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return err
	}
	nowUnix := time.Now().Unix()
	timer := time.Duration(parsedToken.Claims.(*jwt.StandardClaims).ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	return j.dao.RDB.SetEX(ctx, fmt.Sprintf(jwtBlacklistKey, cypher.MD5(parsedToken.Raw)), nowUnix, timer).Err()
}

func (j *JwtRepo) GetSecretKey() string {
	return j.cfg.Secret
}

func (j *JwtRepo) GetIssuer() string {
	return j.cfg.Issuer
}
