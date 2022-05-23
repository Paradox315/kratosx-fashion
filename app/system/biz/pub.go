package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jassue/go-storage/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	api "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/data/model"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/cypher"
	"kratosx-fashion/pkg/xsync"
	"os"
	"path"
	"time"
)

type PublicUsecase struct {
	userRepo     UserRepo
	userRoleRepo UserRoleRepo
	captchaRepo  CaptchaRepo
	jwtRepo      JwtRepo
	disk         storage.Storage
	log          *log.Helper
	rdb          *redis.Client
}

func NewPublicUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, captchaRepo CaptchaRepo, jwtRepo JwtRepo, disk storage.Storage, rdb *redis.Client, logger log.Logger) *PublicUsecase {
	return &PublicUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		captchaRepo:  captchaRepo,
		jwtRepo:      jwtRepo,
		disk:         disk,
		log:          log.NewHelper(log.With(logger, "biz", "public")),
		rdb:          rdb,
	}
}

func (p *PublicUsecase) buildJwtUser(ctx context.Context, upo *model.User) (user User, err error) {
	user = User{
		Id:       upo.ID,
		Username: upo.Username,
		Email:    upo.Email,
		Mobile:   upo.Mobile,
		Nickname: upo.Nickname,
		Roles:    nil,
	}
	urs, err := p.userRoleRepo.SelectAllByUserID(ctx, upo.ID)
	if err != nil {
		err = errors.Wrap(err, "PublicUsecase.buildUserDto.SelectAllByUserID")
		p.log.WithContext(ctx).Error(err)
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, ur.RoleID)
	}
	for _, rid := range rids {
		user.Roles = append(user.Roles, UserRole{Id: rid})
	}
	return
}
func (p *PublicUsecase) Register(ctx context.Context, regInfo RegisterInfo, c Captcha) (err error) {
	if os.Getenv("env") != "dev" {
		if !p.captchaRepo.Verify(ctx, c) {
			err = api.ErrorCaptchaInvalid("验证码错误")
			return
		}
	}
	user := &model.User{
		Username: regInfo.Username,
		Email:    regInfo.Email,
		Mobile:   regInfo.Mobile,
		Password: regInfo.Password,
	}
	user.Password = cypher.BcryptMake(regInfo.Password)
	if err = p.userRepo.Insert(ctx, user); err != nil {
		p.log.WithContext(ctx).Error(err)
		err = kerrors.BadRequest("REGISTER_FAILED", "注册失败,请检查用户名是否已存在")
	}
	if err = p.userRoleRepo.Insert(ctx, &model.UserRole{UserID: user.ID, RoleID: 3}); err != nil {
		p.log.WithContext(ctx).Error(err)
		err = kerrors.BadRequest("REGISTER_FAILED", "注册失败,请检查用户名是否已存在")
	}
	return
}

func (p *PublicUsecase) Login(ctx context.Context, loginSession UserSession, c Captcha) (token *Token, err error) {
	if !p.captchaRepo.Verify(ctx, c) {
		err = api.ErrorCaptchaInvalid("验证码错误")
		return
	}
	upo, err := p.userRepo.SelectByUsername(ctx, loginSession.Username)
	if err != nil || !cypher.BcryptCheck(loginSession.Password, upo.Password) {
		err = api.ErrorUserInvalid("用户名或密码错误")
		return
	}
	if upo.Status == model.UserStatusForbid {
		err = api.ErrorUserInvalid("用户已被禁用")
		return
	}
	user, err := p.buildJwtUser(ctx, upo)
	if err != nil {
		return
	}
	token, err = p.jwtRepo.Create(ctx, user)
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		err = api.ErrorTokenGenerateFail("生成token失败")
		return
	}
	if err = ctxutil.SetUid(ctx, user.Id); err != nil {
		return
	}
	return
}

func (p *PublicUsecase) Refresh(ctx context.Context, refreshToken string) (token *Token, err error) {
	claims, err := p.jwtRepo.ParseToken(ctx, refreshToken)
	if err != nil {
		return
	}
	lock := xsync.Lock("refresh_lock", int64(time.Second*2), p.rdb)
	// 生成新的token
	if lock.Get() {
		defer lock.Release()
		var upo *model.User
		upo, err = p.userRepo.Select(ctx, cast.ToUint(claims.UID))
		if err != nil {
			return
		}
		var user User
		user, err = p.buildJwtUser(ctx, upo)
		if err != nil {
			return
		}
		token, err = p.jwtRepo.Create(ctx, user)
		if err != nil {
			return
		}
		// 将refreshToken注销
		if err = p.jwtRepo.JoinInBlackList(ctx, refreshToken); err != nil {
			return
		}
		// 保存上下文
		err = ctxutil.SetUid(ctx, user.Id)
	}
	return
}

func (p *PublicUsecase) Generate(ctx context.Context) (c Captcha, err error) {
	return p.captchaRepo.Create(ctx)
}

func (p *PublicUsecase) Logout(ctx context.Context) (err error) {
	lock := xsync.Lock("logout_lock", int64(time.Second*2), p.rdb)
	if lock.Get() {
		defer lock.Release()
		var tokens []string
		tokens, err = p.userRepo.SelectTokens(ctx, ctxutil.GetUid(ctx))
		if err != nil {
			return
		}
		for _, token := range tokens {
			if err = p.jwtRepo.JoinInBlackList(ctx, token); err != nil {
				return
			}
		}
	}
	return
}

func (p *PublicUsecase) Upload(ctx context.Context, params UploadInfo) (url string, err error) {
	file, err := params.File.Open()
	if err != nil {
		err = api.ErrorFileOpenFail("文件打开失败")
		p.log.WithContext(ctx).Error(err)
		return
	}
	fileSuffix := path.Ext(params.File.Filename)
	fid, _ := uuid.NewRandom()
	key := fid.String() + fileSuffix
	err = p.disk.Put(key, file, params.File.Size)
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		err = api.ErrorFileUploadFail("文件上传失败")
		return
	}
	url = p.disk.Url(key)
	return
}
