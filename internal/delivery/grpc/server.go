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
	log      *zap.Logger
}

func NewServer(v usecase.VerifyUseCase, s usecase.SignupUseCase, log *zap.Logger) *GRPCServer {
	return &GRPCServer{verifyUC: v, signupUC: s, log: log}
}

func RegisterGRPCServer(gs *grpc.Server, srv *GRPCServer) {
	pb.RegisterAuthServiceServer(gs, srv)
}
