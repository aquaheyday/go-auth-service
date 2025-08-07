package grpc

import (
	"github.com/aquaheyday/go-auth-service/internal/usecase"
	"github.com/aquaheyday/go-auth-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	proto.UnimplementedAuthServiceServer
	verifyUC usecase.VerifyUsecase
	signupUC usecase.SignupUsecase
	log      *zap.Logger
}

func NewServer(v usecase.VerifyUsecase, s usecase.SignupUsecase, log *zap.Logger) *GRPCServer {
	return &GRPCServer{verifyUC: v, signupUC: s, log: log}
}

func RegisterGRPCServer(gs *grpc.Server, srv *GRPCServer) {
	proto.RegisterAuthServiceServer(gs, srv)
}
