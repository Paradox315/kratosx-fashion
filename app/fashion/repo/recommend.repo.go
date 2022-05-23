package repo

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"kratosx-fashion/app/fashion/biz"
	"kratosx-fashion/app/fashion/model"
	"kratosx-fashion/app/system/conf"
)

const (
	ItemNeighbor  = "/api/item/%s/neighbors"
	UserNeighbor  = "/api/user/%s/neighbors"
	Latest        = "/api/latest"
	Popular       = "/api/popular"
	UserRecommend = "/api/recommend/%s"
)

type RecommendRepo struct {
	httpCli *fiber.Agent
	log     *log.Helper
	target  string
}

func NewRecommendRepo(algo *conf.Algorithm, logger log.Logger) biz.RecommendRepo {
	return &RecommendRepo{
		log:    log.NewHelper(log.With(logger, "repo", "recommend")),
		target: HTTP + algo.RecommendAddr,
	}
}

func (r *RecommendRepo) SelectUserNeighbors(ctx context.Context, uid string, limit, offset int) (items []string, err error) {
	r.httpCli = fiber.AcquireAgent()
	req := r.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(r.target+UserNeighbor, uid))
	req.Header.SetMethod(fiber.MethodGet)
	r.httpCli.QueryString(fmt.Sprintf("n=%d&offset=%d", limit, offset))
	if err = r.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "repo.recommend.SelectUserNeighbors")
		r.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := r.httpCli.Bytes()
	if code != fiber.StatusOK || errs != nil {
		r.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &items); err != nil {
		return
	}
	return
}

func (r *RecommendRepo) SelectClothesNeighbors(ctx context.Context, id string, limit, offset int) (items []string, err error) {
	r.httpCli = fiber.AcquireAgent()
	req := r.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(r.target+ItemNeighbor, id))
	req.Header.SetMethod(fiber.MethodGet)
	r.httpCli.QueryString(fmt.Sprintf("n=%d&offset=%d", limit, offset))
	if err = r.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "repo.recommend.SelectClothesNeighbors")
		r.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := r.httpCli.Bytes()
	if code != fiber.StatusOK || errs != nil {
		r.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &items); err != nil {
		return
	}
	return
}

func (r *RecommendRepo) SelectPopular(ctx context.Context, limit, offset int) (items []model.Item, err error) {
	r.httpCli = fiber.AcquireAgent()
	req := r.httpCli.Request()
	req.SetRequestURI(r.target + Popular)
	req.Header.SetMethod(fiber.MethodGet)
	r.httpCli.QueryString(fmt.Sprintf("n=%d&offset=%d", limit, offset))
	if err = r.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "repo.recommend.SelectPopular")
		r.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := r.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		r.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &items); err != nil {
		return
	}
	return
}

func (r *RecommendRepo) SelectLatest(ctx context.Context, limit, offset int) (items []model.Item, err error) {
	r.httpCli = fiber.AcquireAgent()
	req := r.httpCli.Request()
	req.SetRequestURI(r.target + Latest)
	req.Header.SetMethod(fiber.MethodGet)
	r.httpCli.QueryString(fmt.Sprintf("n=%d&offset=%d", limit, offset))
	if err = r.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "repo.recommend.SelectLatest")
		r.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := r.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		r.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &items); err != nil {
		return
	}
	return
}

func (r *RecommendRepo) SelectUserRecommend(ctx context.Context, uid string, limit, offset int) (items []string, err error) {
	r.httpCli = fiber.AcquireAgent()
	req := r.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(r.target+UserRecommend, uid))
	req.Header.SetMethod(fiber.MethodGet)
	r.httpCli.QueryString(fmt.Sprintf("n=%d&offset=%d&write-back-type=read", limit, offset))
	if err = r.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "repo.recommend.SelectUserRecommend")
		r.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := r.httpCli.Bytes()
	if code != fiber.StatusOK || errs != nil {
		r.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &items); err != nil {
		return
	}
	return
}
