package xgrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"
	
	"github.com/jinares/gopkg/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

//GetServerUnaryInvokerLog grpc server call log
func ServerUnaryInvokerLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		
		tlog := xlog.CtxLog(ctx)
		t1 := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(t1)
		springtf := fmt.Sprintf("grpc-server-method:%s; req: %+v; resp: %+v; error: %+v; elapsed: %s", info.FullMethod, req, resp, err, elapsed)
		if err != nil {
			tlog.Error(springtf)
		} else {
			tlog.Info(springtf)
		}
		return resp, err
	}

}

//TLSServerOpt TLSServerOpt
func TLSServerOpt(tlsCert, tlsKey, caCert string) grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		ulog.Fatal(err)
	}
	rawCaCert, err := ioutil.ReadFile(caCert)
	if err != nil {
		ulog.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rawCaCert)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	return grpc.Creds(creds)
}

//ServerUnaryOpentracing rewrite server's interceptor with open tracing
func ServerUnaryOpentracing(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if tracer == nil {
			return handler(ctx, req)
		}
		//从context中取出metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		//从metadata中取出最终数据，并创建出span对象
		spanContext, err := tracer.Extract(opentracing.TextMap, xlog.MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			fmt.Errorf("extract from metadata err %v", err)
		}
		//初始化server 端的span
		serverSpan := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCServer,
		)
		defer serverSpan.Finish()
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		//将带有追踪的context传入应用代码中进行调用
		return handler(ctx, req)
	}
}

//CreateWithUnaryServerChain WithUnaryServerChain
func CreateWithUnaryServerChain(item ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(item...)
}

//CreateWithStreamServerChain  WithStreamServerChain
func CreateWithStreamServerChain(item ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(item...)
}
func AdditionalUnaryServerOptWithBase(item ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	ret := []grpc.UnaryServerInterceptor{
		grpc_validator.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
	}
	return append(ret, item...)
}

//ServerUnaryRecoveryInterceptor RecoveryInterceptor,
func ServerUnaryRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			//debug.PrintStack()
			msg := debug.Stack()
			log:=xlog.CtxLog(ctx)
			log.Errorf("gprc-recover-fail:%s; msg:%s  %v", info.FullMethod, string(msg), e)
			
		}
	}()
	
	return handler(ctx, req)
}