package grpc

import (
	"context"
	"github.com/aquaheyday/go-auth-service/proto"
	"go.uber.org/zap"
)

func (s *GRPCServer) SendVerification(ctx context.Context, req *proto.SendVerificationReq) (*proto.SendVerificationRes, error) {
	if err := s.verifyUC.SendVerification(ctx, req.Email); err != nil {
		s.log.Error("SendVerification failed", zap.Error(err))
		return nil, err
	}
	return &proto.SendVerificationRes{Message: "Verification code sent"}, nil
}

func (s *GRPCServer) VerifyCode(ctx context.Context, req *proto.VerifyCodeReq) (*proto.VerifyCodeRes, error) {
	ok, err := s.verifyUC.VerifyCode(ctx, req.Email, req.Code)
	if err != nil {
		s.log.Error("VerifyCode failed", zap.Error(err))
		return nil, err
	}
	return &proto.VerifyCodeRes{Ok: ok}, nil
}

func (s *GRPCServer) SignUp(ctx context.Context, req *proto.SignUpReq) (*proto.SignUpRes, error) {
	userID, err := s.signupUC.SignUp(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		s.log.Error("SignUp failed", zap.Error(err))
		return nil, err
	}
	return &proto.SignUpRes{UserId: userID}, nil
}
