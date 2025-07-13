package logging

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("[gRPC] → request: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("[gRPC] → error: %v", err)
		}
		return resp, err
	}
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		log.Printf("[gRPC] → %s req=%+v", method, req)

		err := invoker(ctx, method, req, reply, cc, opts...)

		log.Printf("[gRPC] ← %s resp=%+v err=%v took=%s", method, reply, err, time.Since(start))
		return err
	}
}
