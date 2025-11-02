package grpc

import (
	"github.com/aquaheyday/go-auth-service/internal/usecase"
	pb "github.com/aquaheyday/go-auth-service/pkg/pb/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedAuthServiceServer
	verifyUC usecase.VerifyUseCase
	signupUC usecase.SignupUseCase
	loginUC  usecase.LoginUseCase
	log      *zap.Logger
}

func RegisterGRPCServer(gs *grpc.Server, srv *GRPCServer) {
	pb.RegisterAuthServiceServer(gs, srv)
}

func NewGRPCServer(
	logger *zap.Logger,
	verifyUC usecase.VerifyUseCase,
	signupUC usecase.SignupUseCase,
	loginUC usecase.LoginUseCase, // 생성자에 파라미터 추가
) *GRPCServer {
	return &GRPCServer{
		log:      logger,
		verifyUC: verifyUC,
		signupUC: signupUC,
		loginUC:  loginUC, // 필드 초기화
	}
}
