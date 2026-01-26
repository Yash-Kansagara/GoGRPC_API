package interceptors

import (
	"context"
	"log"
	"os"
	"strconv"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var limiter *rate.Limiter

func RateLimitingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if !limiter.Allow() {
		log.Println("Rate limit exceeded")
		return nil, status.Error(codes.ResourceExhausted, "Too many requests")
	}

	return handler(ctx, req)
}

func InitRateLimiter() {
	if limiter == nil {
		limit, err := strconv.Atoi(os.Getenv("GRPC_RATE_LIMIT"))
		if err != nil {
			limit = 10
			log.Println("Failed to read rate limit from env, using default value", limit)
		}

		burst, err := strconv.Atoi(os.Getenv("GRPC_RATE_BURST"))
		if err != nil {
			burst = 10
			log.Println("Failed to read rate limit burst from env, using default value", burst)
		}
		limiter = rate.NewLimiter(rate.Limit(limit), burst)
	}
}
