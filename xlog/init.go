package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	envServiceName = "JAEGER_SERVICE_NAME"

	ctxTraceID  = "ctx.traceid"
	ctxSpanID   = "ctx.spanid"
	ctxParentID = "ctx.parentid"
)

func init() {
	xLOG.SetOutput(os.Stderr)
	xLOG.SetLevel(logrus.TraceLevel)
}
