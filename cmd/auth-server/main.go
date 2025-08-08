// cmd/auth-server/main.go
// 이 파일은 auth-service의 진입점이며, 설정 로드, 의존성 초기화,
// 그리고 gRPC 서버를 구성하고 실행하는 역할을 합니다.
package main

import (
	mailerinfra "github.com/aquaheyday/go-auth-service/internal/infra/mailer/mock"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	grpcdeliv "github.com/aquaheyday/go-auth-service/internal/delivery/grpc"
	"github.com/aquaheyday/go-auth-service/internal/infra/cache"
	"github.com/aquaheyday/go-auth-service/internal/infra/db"
	postgresrepo "github.com/aquaheyday/go-auth-service/internal/repository/postgres"
	redisrepo "github.com/aquaheyday/go-auth-service/internal/repository/redis"
	"github.com/aquaheyday/go-auth-service/internal/usecase"
	"github.com/aquaheyday/go-auth-service/pkg/config"
	"github.com/aquaheyday/go-auth-service/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 설정 파일 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		// 설정 로드 실패 시 즉시 종료
		log.Fatalf("failed to load config: %v", err)
	}

	// 로거(logger) 초기화
	logg := logger.NewLogger(cfg.LogLevel)
	defer logg.Sync() // 애플리케이션 종료 시 로그 버퍼 플러시

	// Postgres 데이터베이스 연결
	postgresDbConn, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		logg.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer postgresDbConn.Close() // 애플리케이션 종료 시 연결 해제

	// Redis 캐시 클라이언트 생성
	rdb := cache.NewRedis(cfg.RedisAddr)
	defer rdb.Close() // 애플리케이션 종료 시 연결 해제

	// 메일러(메일 발송) 인터페이스 초기화 (mock 사용)
	mailSender := mailerinfra.NewSMTPMailer(
		cfg.SMTPHost, cfg.SMTPPort,
		cfg.SMTPUser, cfg.SMTPPass,
	)

	// 레포지토리 및 유스케이스(비즈니스 로직) 구성
	userRepo := postgresrepo.NewUserRepository(postgresDbConn)   // 사용자 저장소
	verificationRepo := redisrepo.NewVerificationRepository(rdb) // 검증 코드 저장소
	// 검증 코드 발송 및 확인 유스케이스
	verifyUC := usecase.NewVerifyUseCase(verificationRepo, mailSender)
	// 회원가입 유스케이스
	signupUC := usecase.NewSignupUseCase(userRepo, verificationRepo)

	// gRPC 서버 리스너 생성
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logg.Fatal("failed to listen", zap.Error(err))
	}

	// gRPC 서버 인스턴스 및 핸들러 등록
	server := grpcdeliv.NewServer(verifyUC, signupUC, logg)
	grpcServer := grpc.NewServer()
	grpcdeliv.RegisterGRPCServer(grpcServer, server)

	// gRPC 리플렉션 서비스 등록
	// 클라이언트에서 동적으로 서비스 정보를 조회할 수 있도록 함
	reflection.Register(grpcServer)

	// gRPC 서버 실행
	logg.Info("gRPC server running", zap.String("port", cfg.GRPCPort))
	if err := grpcServer.Serve(lis); err != nil {
		logg.Fatal("failed to serve grpc", zap.Error(err))
	}
}
