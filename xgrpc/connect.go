package xgrpc

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

var (
	lock     sync.RWMutex
	grpcConn = map[string]*grpc.ClientConn{}
)

func GrpcConn(addr string) (*grpc.ClientConn, error) {
	lock.RLock()
	conn, isok := grpcConn[addr]

	lock.RUnlock()
	if isok {
		return conn, nil
	}
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	serviceName := os.Getenv("JAEGER_SERVICE_NAME")
	opt := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				ClientUnaryOpentracingAndTimeOut(serviceName, 6*time.Second),
			),
		),
	}
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()
	conn, err := grpc.DialContext(ctx, addr, opt...)
	if err != nil {
		ulog.Errorf("grpc-connection-fail:%s", err.Error())
		return nil, err
	}
	grpcConn[addr] = conn
	return conn, nil
}
func Close() {
	lock.Lock()
	defer lock.Unlock()
	for _, conn := range grpcConn {
		conn.Close()
	}
	grpcConn = map[string]*grpc.ClientConn{}
}
