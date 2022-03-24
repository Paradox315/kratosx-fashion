package biz

import (
	"context"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/cypher"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"

	kmw "github.com/go-kratos/kratos/v2/middleware"
	api "kratosx-fashion/api/system/v1"
)

type PublicUsecase struct {
	userRepo    UserRepo
	logRepo     LoginLogRepo
	captchaRepo CaptchaRepo
	jwtSrv      kmw.FiberMiddleware
	log         *log.Helper
}

func NewPublicUsecase(userRepo UserRepo, logRepo LoginLogRepo, captchaRepo CaptchaRepo, jwtSrv kmw.FiberMiddleware, logger log.Logger) *PublicUsecase {
	return &PublicUsecase{
		userRepo:    userRepo,
		logRepo:     logRepo,
		captchaRepo: captchaRepo,
		jwtSrv:      jwtSrv,
		log:         log.NewHelper(logger),
	}
}

func (p *PublicUsecase) Register(ctx context.Context, regInfo RegisterInfo, c Captcha) (uid string, username string, err error) {
	if !p.captchaRepo.Verify(ctx, c.CaptchaId, c.Captcha) {
		err = api.ErrorCaptchaInvalid("验证码错误")
		return
	}
	var user model.User
	if p.userRepo.ExistByUserName(ctx, regInfo.Username) {
		err = api.ErrorUserAlreadyExists("用户名已存在")
		return
	}

	if len(regInfo.Email) > 0 {
		if p.userRepo.ExistByEmail(ctx, regInfo.Email) {
			err = api.ErrorEmailAlreadyExists("邮箱已存在")
			return
		}
	}

	if len(regInfo.Mobile) > 0 {
		if p.userRepo.ExistByMobile(ctx, regInfo.Mobile) {
			err = api.ErrorMobileAlreadyExists("手机号已存在")
			return
		}
	}

	if err = copier.Copy(&user, &regInfo); err != nil {
		return
	}
	user.Password = cypher.BcryptMake(regInfo.Password)
	err = p.userRepo.Insert(ctx, &user)
	username = user.UserName
	uid = user.GetUid()
	return
}

func (p *PublicUsecase) Login(ctx context.Context, loginSession UserSession, c Captcha) (user model.User, err error) {
	if !p.captchaRepo.Verify(ctx, c.CaptchaId, c.Captcha) {
		err = api.ErrorCaptchaInvalid("验证码错误")
		return
	}

	return
}
