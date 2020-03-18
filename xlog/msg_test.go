package xlog

import (
	"fmt"
	"github.com/jinares/gopkg/xtools"
	"os"
	"testing"
)

func TestMsg_Add(t *testing.T) {
	xLOG.SetOutput(os.Stderr)
	xLOG.SetFormatter(&TextLoggerFormatter{})

	msg := NewMsg("out-act").Add(123).BR().Add(map[string]interface{}{
		"msg": xtools.Guid(),
	})
	xLOG.Info(msg.Add("9999"))
	xLOG.Info(msg)
	fmt.Sprint(msg)
}
