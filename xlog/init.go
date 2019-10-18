package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	tlog.SetOutput(os.Stderr)
	tlog.SetLevel(logrus.TraceLevel)
}
