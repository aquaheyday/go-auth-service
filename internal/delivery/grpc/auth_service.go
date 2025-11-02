package grpc

import (
	"context"
	pb "github.com/aquaheyday/go-auth-service/pkg/pb/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	userID, accessToken, refreshToken, err := s.loginUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.log.Error("Login failed", zap.Error(err))
		return nil, err
	}

	return &pb.LoginRes{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *authServiceServer) SendPhoneVerification(ctx context.Context, req *auth.SendPhoneVerificationReq) (*auth.SendPhoneVerificationRes, error) {
	// 전화번호 유효성 검사
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}

	// 유스케이스 호출
	if err := s.verifyUseCase.SendPhoneVerification(ctx, req.PhoneNumber); err != nil {
		s.logger.Error("Failed to send phone verification", zap.Error(err), zap.String("phone", req.PhoneNumber))
		return nil, status.Errorf(codes.Internal, "failed to send verification: %v", err)
	}

	return &auth.SendPhoneVerificationRes{
		Message: "Verification code sent to your phone",
	}, nil
}

func (s *authServiceServer) VerifyPhoneCode(ctx context.Context, req *auth.VerifyPhoneCodeReq) (*auth.VerifyPhoneCodeRes, error) {
	// 입력값 유효성 검사
	if req.PhoneNumber == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number and code are required")
	}

	// 유스케이스 호출
	verified, err := s.verifyUseCase.VerifyPhoneCode(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		s.logger.Error("Failed to verify phone code", zap.Error(err), zap.String("phone", req.PhoneNumber))
		return nil, status.Errorf(codes.Internal, "verification failed: %v", err)
	}

	if !verified {
		return &auth.VerifyPhoneCodeRes{Ok: false}, nil
	}

	return &auth.VerifyPhoneCodeRes{Ok: true}, nil
}
