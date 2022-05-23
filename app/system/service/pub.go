package service

import (
	"context"
	"kratosx-fashion/app/system/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"

	pb "kratosx-fashion/api/system/v1"
	mw "kratosx-fashion/app/system/middleware"
)

type PubService struct {
	pb.UnimplementedPubServer

	uc     *biz.PublicUsecase
	jwtSrv *mw.JWTService
	log    *log.Helper
}

func NewPubService(uc *biz.PublicUsecase, logger log.Logger) *PubService {
	return &PubService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "service", "public")),
	}
}
func (s *PubService) Generate(ctx context.Context, req *pb.EmptyRequest) (*pb.CaptchaReply, error) {
	c, err := s.uc.Generate(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CaptchaReply{
		CaptchaId: c.CaptchaId,
		PicPath:   c.Captcha,
	}, nil
}
func (s *PubService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.EmptyReply, error) {
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
	err := s.uc.Register(ctx, r, c)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
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
	token, err := s.uc.Login(ctx, u, c)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
		TokenType:    token.TokenType,
	}, err
}
func (s *PubService) RefreshToken(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshReply, error) {
	token, err := s.uc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshReply{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
		TokenType:    token.TokenType,
	}, nil
}
func (s *PubService) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.EmptyReply, error) {
	err := s.uc.Logout(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
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
