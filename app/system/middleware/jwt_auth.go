package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"
	"kratosx-fashion/app/system/biz"
	"os"
	"strings"
	"sync"
)

var _ kmw.FiberMiddleware = (*JWTService)(nil)

const (
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	jwtReason string = "JWT_AUTH_ERROR"
)

var (
	ErrMissingJwtToken = errors.Unauthorized(jwtReason, "JWT token is missing")
)

type JWTService struct {
	jwtRepo  biz.JwtRepo
	userRepo biz.UserRepo
	once     sync.Once
	log      *log.Helper
}

func NewJwtService(jwtRepo biz.JwtRepo, userRepo biz.UserRepo, logger log.Logger) *JWTService {
	j := &JWTService{
		jwtRepo:  jwtRepo,
		userRepo: userRepo,
		log:      log.NewHelper(logger),
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
			if !j.userRepo.Verify(ctx, claims.UID) {
				return errors.Unauthorized(jwtReason, "invalid user")
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
