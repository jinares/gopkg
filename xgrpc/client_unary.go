package xgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/jinares/gopkg/xlog"
	"github.com/jinares/gopkg/xtools"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//GetUnaryTimeOut
func ClientUnaryOpentracingAndTimeOut(servername string,timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancelFunc := context.WithTimeout(ctx, timeout)
		defer cancelFunc()
		tracer:=opentracing.GlobalTracer()
		//从context中获取spanContext,如果上层没有开启追踪，则这里新建一个
		//追踪，如果上层已经有了，测创建子span．
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		cliSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
		)
		defer cliSpan.Finish()
		
		//将之前放入context中的metadata数据取出，如果没有则新建一个metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		mdWriter := xlog.MDReaderWriter{md}
		
		//将追踪数据注入到metadata中
		err := tracer.Inject(cliSpan.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			fmt.Errorf("inject to metadata err %v", err)
		}
		//将metadata数据装入context中
		ctx = metadata.NewOutgoingContext(ctx, md)
		//使用带有追踪数据的context进行grpc调用．
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			cliSpan.LogFields(log.String("err", err.Error()))
		}
		return err
	}
}


//GetInvokerLog grpc client call log
func ClientUnaryInvokerLog() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tlog := xlog.CtxLog(ctx)
		t1 := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		elapsed := time.Since(t1)
		springf := fmt.Sprintf(
			"grpc-client-method: %s; req: %+v; reply: %+v;error: %+v; elapsed: %s",
			method, xtools.JSONToStr(req), xtools.JSONToStr(reply), err, elapsed,
		)
		if err != nil {
			tlog.Error(springf)
		} else {
			tlog.Info(springf)
		}
		return err
	}

}
