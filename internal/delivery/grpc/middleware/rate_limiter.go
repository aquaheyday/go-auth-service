// internal/delivery/grpc/middleware/rate_limiter.go

package middleware

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type limiterClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*limiterClient
	limit   rate.Limit
	burst   int
	ttl     time.Duration
}

func NewRateLimiter(rps float64, burst int, ttl time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*limiterClient),
		limit:   rate.Limit(rps),
		burst:   burst,
		ttl:     ttl,
	}
}

func (l *RateLimiter) getClientLimiter(clientIP string) *rate.Limiter {
	client, exists := l.clients[clientIP]
	if !exists {
		client = &limiterClient{
			limiter:  rate.NewLimiter(l.limit, l.burst),
			lastSeen: time.Now(),
		}
		l.clients[clientIP] = client
	}
	client.lastSeen = time.Now()
	return client.limiter
}

func (l *RateLimiter) RateLimiterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 클라이언트 IP 추출
		clientIP := extractClientIP(ctx)

		limiter := l.getClientLimiter(clientIP)
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

// extractClientIP는 여러 방법으로 클라이언트 IP를 추출합니다
func extractClientIP(ctx context.Context) string {
	// 방법 1: peer 정보에서 추출
	if pr, ok := peer.FromContext(ctx); ok {
		if pr.Addr != nil {
			// IP:port 형식에서 IP만 추출
			addr := pr.Addr.String()
			if host, _, err := net.SplitHostPort(addr); err == nil {
				return host
			}
			return addr
		}
	}

	// 방법 2: X-Forwarded-For 헤더에서 추출 (프록시/로드밸런서 환경)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("x-forwarded-for"); len(values) > 0 {
			// 첫 번째 IP가 원래 클라이언트 IP
			ips := strings.Split(values[0], ",")
			return strings.TrimSpace(ips[0])
		}

		// 방법 3: X-Real-IP 헤더에서 추출
		if values := md.Get("x-real-ip"); len(values) > 0 {
			return values[0]
		}
	}

	// 식별 가능한 IP가 없을 경우
	return "unknown"
}
