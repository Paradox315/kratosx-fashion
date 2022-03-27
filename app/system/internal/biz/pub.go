package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/jassue/go-storage/storage"
	"github.com/jinzhu/copier"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/cypher"
	"path"

	api "kratosx-fashion/api/system/v1"
	mw "kratosx-fashion/app/system/internal/middleware"
)

type PublicUsecase struct {
	userRepo    UserRepo
	logRepo     LoginLogRepo
	captchaRepo CaptchaRepo
	jwtSrv      *mw.JWTService
	disk        storage.Storage
	log         *log.Helper
}

func NewPublicUsecase(userRepo UserRepo, logRepo LoginLogRepo, captchaRepo CaptchaRepo, jwtSrv *mw.JWTService, disk storage.Storage, logger log.Logger) *PublicUsecase {
	return &PublicUsecase{
		userRepo:    userRepo,
		logRepo:     logRepo,
		captchaRepo: captchaRepo,
		jwtSrv:      jwtSrv,
		disk:        disk,
		log:         log.NewHelper(log.With(logger, "biz", "public")),
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
	p.log.WithContext(ctx).Error(err)
	username = user.Username
	uid = user.GetUid()
	return
}

func (p *PublicUsecase) Login(ctx context.Context, loginSession UserSession, c Captcha) (token Token, uid uint, err error) {
	if !p.captchaRepo.Verify(ctx, c.CaptchaId, c.Captcha) {
		err = api.ErrorCaptchaInvalid("验证码错误")
		return
	}
	user, err := p.userRepo.SelectPasswordByName(ctx, loginSession.Username)
	if err != nil || !cypher.BcryptCheck(loginSession.Password, user.Password) {
		err = api.ErrorUserInvalid("用户名或密码错误")
		return
	}
	if user.Status != uint8(model.UserStatusForbid) {
		err = api.ErrorUserInvalid("用户已被禁用")
		return
	}
	tokenOut, err := p.jwtSrv.CreateToken(ctx, user)
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	token = Token{
		AccessToken: tokenOut.Token,
		TokenType:   token.TokenType,
		ExpireAt:    token.ExpireAt,
	}
	return
}

func (p *PublicUsecase) Generate(ctx context.Context) (id string, b64s string, err error) {
	return p.captchaRepo.Create(ctx)
}

func (p *PublicUsecase) Logout(ctx context.Context, token string) (err error) {
	return p.jwtSrv.JoinBlackList(ctx, token)
}

func (p *PublicUsecase) Upload(ctx context.Context, params UploadInfo) (url string, err error) {
	file, err := params.File.Open()
	if err != nil {
		err = api.ErrorFileOpenFail("文件打开失败")
		p.log.WithContext(ctx).Error(err)
		return
	}
	fileSuffix := path.Ext(params.File.Filename)
	fid, _ := uuid.NewUUID()
	key := fid.String() + fileSuffix
	err = p.disk.Put(key, file, params.File.Size)
	p.log.WithContext(ctx).Error(err)
	url = p.disk.Url(key)
	return
}
