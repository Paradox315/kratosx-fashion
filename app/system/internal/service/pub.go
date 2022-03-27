package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"kratosx-fashion/app/system/internal/biz"
	"strconv"

	pb "kratosx-fashion/api/system/v1"
)

type PubService struct {
	pb.UnimplementedPubServer

	uc  *biz.PublicUsecase
	log *log.Helper
}

func NewPubService(uc *biz.PublicUsecase, logger log.Logger) *PubService {
	return &PubService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *PubService) Generate(ctx context.Context, req *pb.EmptyRequest) (*pb.CaptchaReply, error) {
	id, b64s, err := s.uc.Generate(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CaptchaReply{
		CaptchaId: id,
		PicPath:   b64s,
	}, nil
}
func (s *PubService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	r := biz.RegisterInfo{
		Email:    req.Email,
		Mobile:   req.Mobile,
		Username: req.Username,
		Password: req.Password,
	}
	c := biz.Captcha{
		Captcha:   req.Captcha,
		CaptchaId: req.CaptchaId,
	}
	uid, username, err := s.uc.Register(ctx, r, c)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{
		UserId:   uid,
		Username: username,
	}, nil
}
func (s *PubService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	u := biz.UserSession{
		Username: req.Username,
		Password: req.Password,
	}
	c := biz.Captcha{
		Captcha:   req.Captcha,
		CaptchaId: req.CaptchaId,
	}
	token, uid, err := s.uc.Login(ctx, u, c)
	return &pb.LoginReply{
		Token: &pb.TokenData{
			AccessToken: token.AccessToken,
			ExpiresAt:   token.ExpireAt,
			TokenType:   token.TokenType,
		},
		UserId:   strconv.Itoa(int(uid)),
		Username: req.Username,
	}, err
}
func (s *PubService) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.EmptyReply, error) {
	c, ok := transport.FromFiberContext(ctx)
	if !ok {
		return nil, errors.InternalServer("CONTEXT PARSE", "find context error")
	}
	err := s.uc.Logout(ctx, c.Locals("token").(string))
	return &pb.EmptyReply{}, err
}
func (s *PubService) RetrievePwd(ctx context.Context, req *pb.RetrieveRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *PubService) UploadFile(ctx context.Context, req *pb.EmptyRequest) (*pb.UploadReply, error) {
	c, ok := transport.FromFiberContext(ctx)
	if !ok {
		return nil, errors.InternalServer("CONTEXT PARSE", "find context error")
	}
	var params biz.UploadInfo
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	params.File = file
	url, err := s.uc.Upload(ctx, params)
	if err != nil {
		return nil, err
	}
	return &pb.UploadReply{
		Url: url,
	}, nil
}
