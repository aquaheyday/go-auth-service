package grpc

import (
	"context"
	pb "github.com/aquaheyday/go-auth-service/pkg/pb/auth"
	"go.uber.org/zap"
)

func (s *GRPCServer) SendVerification(ctx context.Context, req *pb.SendVerificationReq) (*pb.SendVerificationRes, error) {
	if err := s.verifyUC.SendVerification(ctx, req.Email); err != nil {
		s.log.Error("SendVerification failed", zap.Error(err))
		return nil, err
	}
	return &pb.SendVerificationRes{Message: "Verification code sent"}, nil
}

func (s *GRPCServer) VerifyCode(ctx context.Context, req *pb.VerifyCodeReq) (*pb.VerifyCodeRes, error) {
	ok, err := s.verifyUC.VerifyCode(ctx, req.Email, req.Code)
	if err != nil {
		s.log.Error("VerifyCode failed", zap.Error(err))
		return nil, err
	}
	return &pb.VerifyCodeRes{Ok: ok}, nil
}

func (s *GRPCServer) SignUp(ctx context.Context, req *pb.SignUpReq) (*pb.SignUpRes, error) {
	userID, err := s.signupUC.SignUp(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		s.log.Error("SignUp failed", zap.Error(err))
		return nil, err
	}
	return &pb.SignUpRes{UserId: userID}, nil
}
