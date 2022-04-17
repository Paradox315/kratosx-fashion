package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
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
	uid := ctxutil.GetUid(ctx)
	user, err := s.uc.Get(ctx, cast.ToUint(uid))
	if err != nil {
		return nil, err
	}
	resp := &pb.UserState{}
	_ = copier.Copy(&resp, &user)
	resp.RegisterDate = user.CreatedAt
	for _, role := range user.Roles {
		resp.Roles = append(resp.Roles, role.Id)
	}
	return resp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.Save(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.Edit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) (*pb.IDReply, error) {
	err := s.uc.EditPassword(ctx, req.OldPassword, req.NewPassword, req.ConfirmPassword, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: req.Id,
	}, nil
}
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.StatusRequest) (*pb.IDReply, error) {
	err := s.uc.EditStatus(ctx, cast.ToUint(req.Id), model.UserStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: req.Id,
	}, nil
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
	limit, offset := pagination.Parse(req.PageNum, req.PageSize)
	var opt *biz.SQLOption
	if len(req.Query) > 0 {
		where, order, args := option.Parse(req.Query...)
		opt = &biz.SQLOption{
			Where: where,
			Order: order,
			Args:  args,
		}
	}
	users, total, err := s.uc.Search(ctx, limit, offset, opt)
	if err != nil {
		return
	}
	list = &pb.ListUserReply{
		Total: uint32(total),
	}
	for _, user := range users {
		uReply := &pb.UserReply{}
		_ = copier.Copy(&uReply, &user)
		list.List = append(list.List, uReply)
	}
	return
}
func (s *UserService) ListLoginLog(ctx context.Context, req *pb.ListRequest) (list *pb.ListLoginLogReply, err error) {
	limit, offset := pagination.Parse(req.PageNum, req.PageSize)
	uid := ctxutil.GetUid(ctx)
	return s.uc.LogPage(ctx, cast.ToUint64(uid), limit, offset)
}
