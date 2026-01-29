package interceptors

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func ResponseTimeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()

	// logging response time, could be reported to some monitoring tool
	fmt.Println(info.FullMethod, end.Sub(start))
	return resp, err
}
