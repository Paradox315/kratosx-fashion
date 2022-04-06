package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
	"github.com/jassue/go-storage/storage"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/cypher"
	"path"
	"strconv"

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
func (p *PublicUsecase) buildUserDto(ctx context.Context, upo *model.User) (user User, err error) {
	if err = copier.Copy(&user, &upo); err != nil {
		return
	}
	urs, err := p.userRoleRepo.SelectAllByUserID(ctx, uint64(upo.ID))
	if err != nil {
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, uint(ur.RoleID))
	}
	for _, rid := range rids {
		user.UserRoles = append(user.UserRoles, UserRole{ID: strconv.FormatUint(uint64(rid), 10)})
	}
	user.Id = cast.ToString(upo.ID)
	user.CreatedAt = upo.CreatedAt.Format(timeFormat)
	user.UpdatedAt = upo.UpdatedAt.Format(timeFormat)
	user.Gender = upo.Gender.String()
	return
}
func (p *PublicUsecase) Register(ctx context.Context, regInfo RegisterInfo, c Captcha) (uid string, username string, err error) {
	if !p.captchaRepo.Verify(ctx, c) {
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
	uid = cast.ToString(user.ID)
	return
}

func (p *PublicUsecase) Login(ctx context.Context, loginSession UserSession, c Captcha) (token *Token, uid string, err error) {
	//if !p.captchaRepo.Verify(ctx, c) {
	//	err = api.ErrorCaptchaInvalid("验证码错误")
	//	return
	//}
	upo, err := p.userRepo.SelectPasswordByName(ctx, loginSession.Username)
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
	fc, ok := transport.FromFiberContext(ctx)
	if !ok {
		err = errors.InternalServer("CONTEXT PARSE", "find context error")
		return
	}
	loc, err := p.logRepo.SelectLocation(ctx, fc.IP())
	if err != nil {
		return
	}
	agent, err := p.logRepo.SelectAgent(ctx, fc.Get("User-Agent"))
	if err != nil {
		return
	}
	bytes, err := encoding.GetCodec("json").Marshal(loc)
	if err != nil {
		return
	}
	if err = p.logRepo.Insert(ctx, &model.LoginLog{
		UserID:     uint64(upo.ID),
		Ip:         fc.IP(),
		Location:   string(bytes),
		LoginType:  model.LoginType_Login,
		Agent:      fc.Get("User-Agent"),
		OS:         agent.OS,
		Device:     agent.Device,
		DeviceType: agent.DeviceType,
	}); err != nil {
		return
	}

	return
}

func (p *PublicUsecase) Generate(ctx context.Context) (c Captcha, err error) {
	return p.captchaRepo.Create(ctx)
}

func (p *PublicUsecase) Logout(ctx context.Context, token string) (err error) {
	fc, ok := transport.FromFiberContext(ctx)
	if !ok {
		err = errors.InternalServer("CONTEXT PARSE", "find context error")
		return
	}
	loc, err := p.logRepo.SelectLocation(ctx, fc.IP())
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	agent, err := p.logRepo.SelectAgent(ctx, fc.Get("User-Agent"))
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	bytes, err := encoding.GetCodec("json").Marshal(loc)
	if err != nil {
		p.log.WithContext(ctx).Error(err)
		return
	}
	uid := fc.Locals("user_id")
	if err = p.logRepo.Insert(ctx, &model.LoginLog{
		UserID:     cast.ToUint64(uid),
		Ip:         fc.IP(),
		Location:   string(bytes),
		LoginType:  model.LoginType_Logout,
		Agent:      fc.Get("User-Agent"),
		OS:         agent.OS,
		Device:     agent.Device,
		DeviceType: agent.DeviceType,
	}); err != nil {
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
	p.log.WithContext(ctx).Error(err)
	url = p.disk.Url(key)
	return
}
