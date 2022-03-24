package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mojocn/base64Captcha"
	"kratosx-fashion/app/system/internal/biz"
)

type CaptchaRepo struct {
	store base64Captcha.Store
	log   *log.Helper
}

func NewCaptchaRepo(logger log.Logger) biz.CaptchaRepo {
	return &CaptchaRepo{
		store: base64Captcha.DefaultMemStore,
		log:   log.NewHelper(logger),
	}
}

func (c CaptchaRepo) Create(ctx context.Context) (id string, b64s string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, c.store)
	return cp.Generate()
}

func (c CaptchaRepo) Verify(ctx context.Context, id string, captcha string) bool {
	return c.store.Verify(id, captcha, true)
}
