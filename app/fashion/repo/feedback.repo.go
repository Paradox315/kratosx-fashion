package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"kratosx-fashion/app/fashion/biz"
	"kratosx-fashion/app/fashion/model"
	"kratosx-fashion/app/system/conf"
)

const FeedbackApi = "/api/feedback"

type FeedbackRepo struct {
	httpCli *fiber.Agent
	log     *log.Helper
	target  string
}

func NewFeedbackRepo(algo *conf.Algorithm, logger log.Logger) biz.FeedbackRepo {
	return &FeedbackRepo{
		log:    log.NewHelper(log.With(logger, "repo", "feedback")),
		target: HTTP + algo.RecommendAddr,
	}
}

func (f *FeedbackRepo) Insert(ctx context.Context, feedback *model.Feedback) error {
	f.httpCli = fiber.AcquireAgent()
	req := f.httpCli.Request()
	req.SetRequestURI(f.target + FeedbackApi)
	req.Header.SetMethod(fiber.MethodPut)
	f.httpCli.JSON([]*model.Feedback{feedback})
	if err := f.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "feedback.Insert")
		f.log.WithContext(ctx).Error(err)
		return err
	}
	code, _, errs := f.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		f.log.WithContext(ctx).Errorf("code %d,errs %+v", code, errs)
		return errors.New("http response error")
	}
	return nil
}
