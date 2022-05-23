package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mojocn/base64Captcha"
	"github.com/pkg/errors"
	"kratosx-fashion/app/system/biz"
	"os"
)

type CaptchaRepo struct {
	store base64Captcha.Store
	log   *log.Helper
}

func NewCaptchaRepo(logger log.Logger) biz.CaptchaRepo {
	return &CaptchaRepo{
		store: base64Captcha.DefaultMemStore,
		log:   log.NewHelper(log.With(logger, "repo", "captcha")),
	}
}

func (c CaptchaRepo) Create(ctx context.Context) (captcha biz.Captcha, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, c.store)
	id, b64s, err := cp.Generate()
	if err != nil {
		err = errors.Wrap(err, "create captcha error")
		c.log.WithContext(ctx).Error("generate captcha error", err)
		return
	}
	captcha.CaptchaId = id
	captcha.Captcha = b64s
	return
}

func (c CaptchaRepo) Verify(ctx context.Context, captcha biz.Captcha) bool {
	if os.Getenv("env") == "dev" {
		return true
	}
	return c.store.Verify(captcha.CaptchaId, captcha.Captcha, true)
}
