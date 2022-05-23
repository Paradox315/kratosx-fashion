package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/biz"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/option"
	"kratosx-fashion/pkg/pagination"
	"kratosx-fashion/pkg/xcast"
	"strings"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "service", "user")),
	}
}
func (s *UserService) InitUserInfo(ctx context.Context, req *pb.EmptyRequest) (*pb.UserState, error) {
	return s.uc.Init(ctx)
}
func (s *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.Save(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: uint64(id),
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.EmptyReply, error) {
	if err := s.uc.Edit(ctx, req); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) (*pb.EmptyReply, error) {
	err := s.uc.EditPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) ResetPassword(ctx context.Context, req *pb.IDRequest) (*pb.EmptyReply, error) {
	err := s.uc.ResetPassword(ctx, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	if err := s.uc.Remove(ctx, xcast.ToUintSlice[string](strings.Split(req.Ids, ","))); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.IDRequest) (*pb.UserReply, error) {
	user, err := s.uc.Get(ctx, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	reply := &pb.UserReply{}
	if err = copier.Copy(&reply, &user); err != nil {
		return nil, err
	}
	return reply, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListSearchRequest) (list *pb.ListUserReply, err error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	var opt *biz.SQLOption
	if len(req.Query) > 0 {
		where, order, args := option.Parse(req.Query...)
		opt = &biz.SQLOption{
			Where: where,
			Order: order,
			Args:  args,
		}
	}
	return s.uc.Search(ctx, limit, offset, opt)
}
func (s *UserService) ListUserLog(ctx context.Context, req *pb.ListRequest) (list *pb.ListUserLogReply, err error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	uid := ctxutil.GetUid(ctx)
	return s.uc.LogPage(ctx, uid, limit, offset)
}
