package biz

import (
	"context"
	"github.com/pkg/errors"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/cypher"
	"os"
	"path"
	"strconv"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
	"github.com/jassue/go-storage/storage"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	api "kratosx-fashion/api/system/v1"
)

type PublicUsecase struct {
	userRepo     UserRepo
	userRoleRepo UserRoleRepo
	logRepo      LoginLogRepo
	captchaRepo  CaptchaRepo
	jwtRepo      JwtRepo
	disk         storage.Storage
	log          *log.Helper
}

func NewPublicUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, logRepo LoginLogRepo, captchaRepo CaptchaRepo, jwtRepo JwtRepo, disk storage.Storage, logger log.Logger) *PublicUsecase {
	return &PublicUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		logRepo:      logRepo,
		captchaRepo:  captchaRepo,
		jwtRepo:      jwtRepo,
		disk:         disk,
		log:          log.NewHelper(log.With(logger, "biz", "public")),
	}
}
func (p *PublicUsecase) buildLog(ctx context.Context, uid uint64, typ model.LoginType) (log *model.LoginLog, err error) {
	fc, ok := transport.FromFiberContext(ctx)
	if !ok {
		err = kerrors.InternalServer("CONTEXT PARSE", "find context error")
		return
	}
	agent, err := p.logRepo.SelectAgent(ctx, fc.Get("User-Agent"))
	if err != nil {
		err = errors.Wrap(err, "PublicUsecase.Login.SelectAgent")
		p.log.WithContext(ctx).Error(err)
		return
	}

	loc, err := p.logRepo.SelectLocation(ctx, fc.IP())
	if err != nil {
		err = errors.Wrap(err, "PublicUsecase.Login.SelectLocation")
		p.log.WithContext(ctx).Error(err)
		return
	}

	bytes, err := encoding.GetCodec("json").Marshal(&loc)
	if err != nil {
		return
	}
	log = &model.LoginLog{
		UserID:     uid,
		Ip:         fc.IP(),
		Location:   string(bytes),
		LoginType:  typ,
		Agent:      agent.Name,
		OS:         agent.OS,
		Device:     agent.Device,
		DeviceType: agent.DeviceType,
	}
	return
}
func (p *PublicUsecase) buildUserDto(ctx context.Context, upo *model.User) (user User, err error) {
	_ = copier.Copy(&user, &upo)
	urs, err := p.userRoleRepo.SelectAllByUserID(ctx, uint64(upo.ID))
	if err != nil {
		err = errors.Wrap(err, "PublicUsecase.buildUserDto.SelectAllByUserID")
		p.log.WithContext(ctx).Error(err)
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, uint(ur.RoleID))
	}
	for _, rid := range rids {
		user.Roles = append(user.Roles, UserRole{Id: strconv.FormatUint(uint64(rid), 10)})
	}
	user.Id = cast.ToString(upo.ID)
	user.CreatedAt = upo.CreatedAt.Format(timeFormat)
	user.UpdatedAt = upo.UpdatedAt.Format(timeFormat)
	user.Gender = upo.Gender.String()
	return
}
func (p *PublicUsecase) Register(ctx context.Context, regInfo RegisterInfo, c Captcha) (uid string, username string, err error) {
	if os.Getenv("env") != "dev" {
		if !p.captchaRepo.Verify(ctx, c) {
			err = api.ErrorCaptchaInvalid("验证码错误")
			return
		}
	}
	var (
		user model.User
		cnt  int64
	)
	if len(regInfo.Username) > 0 {
		cnt, err = p.userRepo.ExistByUsername(ctx, regInfo.Username)
		if err != nil {
			return
		}
		if cnt > 0 {
			err = api.ErrorUserAlreadyExists("用户名已存在")
			return
		}
	}

	if len(regInfo.Email) > 0 {
		cnt, err = p.userRepo.ExistByEmail(ctx, regInfo.Email)
		if err != nil {
			return
		}
		if cnt > 0 {
			err = api.ErrorEmailAlreadyExists("邮箱已存在")
			return
		}
	}

	if len(regInfo.Mobile) > 0 {
		cnt, err = p.userRepo.ExistByMobile(ctx, regInfo.Mobile)
		if err != nil {
			return
		}
		if cnt > 0 {
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
	uid = cast.ToString(user.ID)
	return
}

func (p *PublicUsecase) Login(ctx context.Context, loginSession UserSession, c Captcha) (token *Token, uid string, err error) {
	if os.Getenv("env") != "dev" {
		if !p.captchaRepo.Verify(ctx, c) {
			err = api.ErrorCaptchaInvalid("验证码错误")
			return
		}
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
	user, err := p.buildUserDto(ctx, upo)
	if err != nil {
		return
	}
	token, err = p.jwtRepo.Create(ctx, user)
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	loginLog, err := p.buildLog(ctx, uint64(upo.ID), model.LoginType_Login)
	if err != nil {
		return
	}
	if err = p.logRepo.Insert(ctx, loginLog); err != nil {
		return
	}

	return
}

func (p *PublicUsecase) Generate(ctx context.Context) (c Captcha, err error) {
	return p.captchaRepo.Create(ctx)
}

func (p *PublicUsecase) Logout(ctx context.Context, token string) (err error) {
	uid := ctxutil.GetUid(ctx)
	loginLog, err := p.buildLog(ctx, cast.ToUint64(uid), model.LoginType_Logout)
	if err != nil {
		return
	}
	if err = p.logRepo.Insert(ctx, loginLog); err != nil {
		return
	}
	return p.jwtRepo.JoinInBlackList(ctx, token)
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
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	url = p.disk.Url(key)
	return
}
