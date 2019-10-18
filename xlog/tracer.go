package xlog

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc/metadata"
)

//XTracer XTracer
var (
	tracer      opentracing.Tracer
	closer      io.Closer
	serviceName = "default"
	tlog        = logrus.New()
)

//Close Close
func Close() error {
	if closer == nil {
		return errors.New("没有初始化")
	}
	err := closer.Close()
	tracer = nil
	closer = nil
	return err
}

//GetTracer GetTracer
func GetTracer() (opentracing.Tracer, error) {
	if tracer == nil {
		return nil, errors.New("init err")
	}

	return tracer, nil

}

//InitByTracerEnv，所有配置均在k8s环境变量中传入
func InitByTracerEnv() error {
	if tracer != nil {
		return nil
	}
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	tracer, closer, err = cfg.NewTracer()
	if err != nil {
		InitTracer(cfg.ServiceName, "127.0.0.1:6831", 10000)
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	serviceName = cfg.ServiceName

	return nil
}

//InitTracer InitTracer
func InitTracer(servicename, dns string, batch int) error {
	if tracer != nil {
		return nil
	}
	udp, _ := jaeger.NewUDPTransport(dns, batch)
	tracer, closer = jaeger.NewTracer(
		servicename,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(udp),
	)
	serviceName = servicename
	opentracing.SetGlobalTracer(tracer)
	return nil
}

type MDReaderWriter struct {
	metadata.MD
}

func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}
func GetLog() *logrus.Logger {

	return tlog
}

//ContextLog ContextLog
func CtxLog(ctx context.Context) *logrus.Entry {

	traceID, spanID, parentSpanID := "", "", ""
	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			sc := span.Context().(jaeger.SpanContext)
			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
			parentSpanID = sc.ParentID().String()
		}
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	mr := MDReaderWriter{md}
	data := map[string]string{}
	mr.ForeachKey(func(key, val string) error {
		data[key] = val
		return nil
	})

	nlog := tlog.WithFields(logrus.Fields{
		"ctx.traceid": traceID,
		"ctx.spanid":  spanID,
		"ctx.pid":     parentSpanID,
	})
	return nlog
}
