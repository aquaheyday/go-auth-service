package main

import (
	"log"
	"net"

	grpcdeliv "github.com/aquaheyday/go-auth-service/internal/delivery/grpc"
	"github.com/aquaheyday/go-auth-service/internal/infra/cache"
	"github.com/aquaheyday/go-auth-service/internal/infra/db"
	mailerinfra "github.com/aquaheyday/go-auth-service/internal/infra/mailer"
	postgresrepo "github.com/aquaheyday/go-auth-service/internal/repository/postgres"
	redisrepo "github.com/aquaheyday/go-auth-service/internal/repository/redis"
	"github.com/aquaheyday/go-auth-service/internal/usecase"
	"github.com/aquaheyday/go-auth-service/pkg/config"
	"github.com/aquaheyday/go-auth-service/pkg/logger"
	"github.com/aquaheyday/go-auth-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 로거 초기화
	logg := logger.NewLogger(cfg.LogLevel)
	defer logg.Sync()

	// Postgres 연결
	dbConn, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		logg.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer dbConn.Close()

	// Redis 연결
	rdb := cache.NewRedis(cfg.RedisAddr)
	defer rdb.Close()

	// 메일러 초기화
	mailSender := mailerinfra.NewSMTPMailer(
		cfg.SMTPHost, cfg.SMTPPort,
		cfg.SMTPUser, cfg.SMTPPass,
	)

	// 레포지토리 및 유스케이스
	userRepo := postgresrepo.NewUserRepository(dbConn)
	verificationRepo := redisrepo.NewVerificationRepository(rdb)
	verifyUC := usecase.NewVerifyUsecase(verificationRepo, mailSender)
	signupUC := usecase.NewSignupUsecase(userRepo, verificationRepo)

	// gRPC 서버 구동
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logg.Fatal("failed to listen", zap.Error(err))
	}
	server := grpcdeliv.NewServer(verifyUC, signupUC, logg)
	grpcServer := grpc.NewServer()
	grpcdeliv.RegisterGRPCServer(grpcServer, server)

	logg.Info("gRPC server running", zap.String("port", cfg.GRPCPort))
	if err := grpcServer.Serve(lis); err != nil {
		logg.Fatal("failed to serve grpc", zap.Error(err))
	}
}
