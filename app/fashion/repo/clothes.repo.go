package repo

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	v1 "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/biz"
	"kratosx-fashion/app/fashion/model"
	"kratosx-fashion/app/system/conf"
	"sync"
)

const (
	InsertAPI = "/api/item"
	SelectAPI = "/api/item/%s"
	UpdateAPI = "/api/item/%s"
	DeleteAPI = "/api/item/%s"
)

type ClothesRepo struct {
	httpCli *fiber.Agent
	grpcCli v1.TryOnClient
	log     *log.Helper
	target  string
}

func (c *ClothesRepo) getItem(ctx context.Context, id string) (item *model.Clothes, err error) {
	cli := fiber.AcquireAgent()
	req := cli.Request()
	req.SetRequestURI(fmt.Sprintf(c.target+SelectAPI, id))
	req.Header.SetMethod(fiber.MethodGet)
	if err = cli.Parse(); err != nil {
		err = errors.Wrap(err, "clothesRepo.Select")
		c.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := cli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		c.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	if err = codec.Unmarshal(bytes, &item); err != nil {
		c.log.WithContext(ctx).Error(err)
		return
	}
	return
}

func NewClothesRepo(algo *conf.Algorithm, logger log.Logger) biz.ClothesRepo {
	repo := &ClothesRepo{
		log:    log.NewHelper(log.With(logger, "repo", "clothes")),
		target: HTTP + algo.RecommendAddr,
	}
	conn, err := grpc.Dial(algo.TryonAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	repo.grpcCli = v1.NewTryOnClient(conn)
	return repo
}

func (c *ClothesRepo) Insert(ctx context.Context, clothes *model.Clothes) (err error) {
	c.httpCli = fiber.AcquireAgent()
	req := c.httpCli.Request()
	req.SetRequestURI(c.target + InsertAPI)
	req.Header.SetMethod(fiber.MethodPost)
	c.httpCli.JSON(clothes)
	if err = c.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "clothesRepo.Insert")
		c.log.WithContext(ctx).Error(err)
		return
	}
	code, _, errs := c.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		c.log.WithContext(ctx).Error(errs)
		return errors.New("http response error")
	}
	return
}

func (c *ClothesRepo) Select(ctx context.Context, id string) (clothes *model.Clothes, err error) {
	c.httpCli = fiber.AcquireAgent()
	req := c.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(c.target+SelectAPI, id))
	req.Header.SetMethod(fiber.MethodGet)
	if err = c.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "clothesRepo.Select")
		c.log.WithContext(ctx).Error(err)
		return
	}
	code, bytes, errs := c.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		c.log.WithContext(ctx).Error(errs)
		err = errors.New("http response error")
		return
	}
	clothes = &model.Clothes{}
	if err = codec.Unmarshal(bytes, clothes); err != nil {
		return
	}
	return
}

func (c *ClothesRepo) SelectByIDs(ctx context.Context, ids []string) (items []*model.Clothes, err error) {
	eg, ctx := errgroup.WithContext(ctx)
	lock := &sync.RWMutex{}
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				c.log.WithContext(ctx).Info("clothesRepo.SelectByIDs context done")
				return nil
			default:
				var item *model.Clothes
				item, err = c.getItem(ctx, id)
				if err != nil {
					return err
				}
				if lock.TryLock() {
					items = append(items, item)
					lock.Unlock()
				} else {
					items = append(items, item)
				}
				return nil
			}
		})
	}
	if err = eg.Wait(); err != nil {
		return
	}
	return
}

func (c *ClothesRepo) Update(ctx context.Context, clothes *model.Clothes) (err error) {
	c.httpCli = fiber.AcquireAgent()
	req := c.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(c.target+UpdateAPI, clothes.ItemId))
	req.Header.SetMethod(fiber.MethodPatch)
	c.httpCli.JSON(clothes)
	if err = c.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "clothesRepo.Update")
		c.log.WithContext(ctx).Error(err)
		return
	}
	code, _, errs := c.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		c.log.WithContext(ctx).Error(errs)
		return errors.New("http response error")
	}
	return
}

func (c *ClothesRepo) Delete(ctx context.Context, id string) (err error) {
	c.httpCli = fiber.AcquireAgent()
	req := c.httpCli.Request()
	req.SetRequestURI(fmt.Sprintf(c.target+DeleteAPI, id))
	req.Header.SetMethod(fiber.MethodDelete)
	if err = c.httpCli.Parse(); err != nil {
		err = errors.Wrap(err, "clothesRepo.Delete")
		c.log.WithContext(ctx).Error(err)
		return
	}
	code, _, errs := c.httpCli.Bytes()
	if len(errs) != 0 || code != fiber.StatusOK {
		c.log.WithContext(ctx).Error(errs)
		return errors.New("http response error")
	}
	return
}

func (c *ClothesRepo) TryOn(ctx context.Context, req *v1.TryOnRequest) (resp *v1.TryOnReply, err error) {
	resp, err = c.grpcCli.TryOnClothes(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "clothesRepo.TryOn")
		c.log.WithContext(ctx).Error(err)
		return
	}
	return resp, nil
}
