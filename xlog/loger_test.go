package xlog

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestCtxLog(t *testing.T) {
	ctx := context.Background()
	xLOG.SetOutput(os.Stderr)
	xLOG.SetFormatter(&logrus.JSONFormatter{})
	CtxLog(ctx).Info("default")
	CtxLog(ctx).Info("default")
	xLOG.SetFormatter(&JSONLoggerFormatter{})
	CtxLog(ctx).Info("json")
	CtxLog(ctx).Info("json")
	xLOG.SetFormatter(&TextLoggerFormatter{})
	CtxLog(ctx).Info("text")
	CtxLog(ctx).Info("text", "test")

}

//
//// 性能测试
//func BenchmarkDefautText(b *testing.B) {
//	ctx := context.Background()
//	tlog.SetOutput(os.Stderr)
//	tlog.SetLevel(logrus.TraceLevel)
//	tlog.SetFormatter(&logrus.TextFormatter{})
//	tlog.SetOutput(&StdWriter{})
//
//	// b.N会根据函数的运行时间取一个合适的值
//	for i := 0; i < b.N; i++ {
//		CtxLog(ctx).Info("default")
//	}
//}

// 性能测试
func BenchmarkMyText(b *testing.B) {
	ctx := context.Background()
	xLOG.SetOutput(os.Stderr)
	xLOG.SetFormatter(&TextLoggerFormatter{})
	xLOG.SetLevel(logrus.TraceLevel)
	//tlog.SetOutput(&StdWriter{})
	// b.N会根据函数的运行时间取一个合适的值
	for i := 0; i < b.N; i++ {
		CtxLog(ctx).Info("default", 1, 2, 3, map[string]string{"a": "aa"})
	}
}
