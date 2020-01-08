package xlog

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/metadata"
)

//XTracer XTracer
var (
	xLOG = logrus.New()
)

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

	return xLOG
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
	//md, ok := metadata.FromIncomingContext(ctx)
	//if !ok {
	//	md = metadata.New(nil)
	//}
	//mr := MDReaderWriter{md}
	//data := map[string]string{}
	//mr.ForeachKey(func(key, val string) error {
	//	data[key] = val
	//	return nil
	//})

	nlog := xLOG.WithFields(logrus.Fields{
		ctxTraceID:  traceID,
		ctxSpanID:   spanID,
		ctxParentID: parentSpanID,
	})
	return nlog
}
