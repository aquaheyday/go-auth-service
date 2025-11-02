// cmd/auth-server/main.go
package main

import (
	"github.com/aquaheyday/go-auth-service/internal/infra/mailer/mock"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"

	grpcdeliv "github.com/aquaheyday/go-auth-service/internal/delivery/grpc"
	"github.com/aquaheyday/go-auth-service/internal/delivery/grpc/middleware" // 미들웨어 패키지 추가
	"github.com/aquaheyday/go-auth-service/internal/infra/cache"
	"github.com/aquaheyday/go-auth-service/internal/infra/db"
	postgresrepo "github.com/aquaheyday/go-auth-service/internal/repository/postgres"
	redisrepo "github.com/aquaheyday/go-auth-service/internal/repository/redis"
	"github.com/aquaheyday/go-auth-service/internal/usecase"
	"github.com/aquaheyday/go-auth-service/pkg/config"
	"github.com/aquaheyday/go-auth-service/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware" // 미들웨어 체인 패키지 추가
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

	// LogLevel이 비어있는 경우 기본값 설정
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info" // 기본값으로 info 레벨 설정
	}

	// 로거(logger) 초기화
	logg := logger.NewLogger(cfg.LogLevel)
	if logg == nil {
		log.Fatal("failed to initialize logger")
	}
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

	// 메일러(메일 발송) 인터페이스 초기화
	/*mailSender := mailerSG.NewSendGridMailer(
		cfg.SendGridAPIKey,
		cfg.SendGridFromEmail,
		cfg.SendGridFromName,
		cfg.SendGridSandbox, // prod에선 false 권장
	)*/

	mailSender := mock.NewSMTPMailer("localhost", 1025, "user", "pass")

	// 레포지토리 및 유스케이스(비즈니스 로직) 구성
	userRepo := postgresrepo.NewUserRepository(postgresDbConn)   // 사용자 저장소
	verificationRepo := redisrepo.NewVerificationRepository(rdb) // 검증 코드 저장소
	// 검증 코드 발송 및 확인 유스케이스
	verifyUC := usecase.NewVerifyUseCase(verificationRepo, mailSender)
	// 회원가입 유스케이스
	signupUC := usecase.NewSignupUseCase(userRepo, verificationRepo)

	// 토큰 레포지토리 생성 및 로그인 유스케이스 추가
	tokenRepo := redisrepo.NewTokenRepository(rdb)          // 토큰 저장소 추가
	loginUC := usecase.NewLoginUseCase(userRepo, tokenRepo) // 로그인 유스케이스 추가

	// gRPC 서버 리스너 생성
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logg.Fatal("failed to listen", zap.Error(err))
	}

	// 미들웨어 설정
	// 1. Prometheus를 사용하여 gRPC 서버 메트릭 수집 미들웨어 생성
	metricsInterceptor := middleware.MetricsInterceptor()

	// 2. 속도 제한 미들웨어 생성 (초당 100 요청, 버스트 200, 1시간 TTL)
	rateLimiter := middleware.NewRateLimiter(100, 200, 1*time.Hour)

	// gRPC 서버 인스턴스 및 핸들러 등록 - 미들웨어 체인 적용
	server := grpcdeliv.NewGRPCServer(logg, verifyUC, signupUC, loginUC)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				metricsInterceptor,                   // 메트릭 미들웨어
				rateLimiter.RateLimiterInterceptor(), // 속도 제한 미들웨어
			),
		),
	)
	grpcdeliv.RegisterGRPCServer(grpcServer, server)

	// gRPC 리플렉션 서비스 등록
	// 클라이언트에서 동적으로 서비스 정보를 조회할 수 있도록 함
	reflection.Register(grpcServer)

	// 메트릭 HTTP 서버 시작
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			logg.Error("Failed to start metrics server", zap.Error(err))
		}
	}()

	// gRPC 서버 실행
	logg.Info("gRPC server running", zap.String("port", cfg.GRPCPort))
	if err := grpcServer.Serve(lis); err != nil {
		logg.Fatal("failed to serve grpc", zap.Error(err))
	}
}
