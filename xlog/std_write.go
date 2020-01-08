package xlog

import (
	"io/ioutil"
	"os"
)

type (
	StdWriter struct {
	}
)

func (h *StdWriter) Write(p []byte) (int, error) {
	err := ioutil.WriteFile("/dev/stderr", p, os.ModeAppend)
	if err != nil {
		return 0, err
	}
	//return os.Stderr.Write(p)
	return len(p), nil
}
