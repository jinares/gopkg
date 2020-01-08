package xlog

import (
	"bytes"
	"fmt"
	"github.com/jinares/gopkg/xtools"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type JSONLoggerFormatter struct {
}

/*

FormatDefault = "[${serviceName},${traceID},${spanID},${parentSpanID}][${pid}]"

timestamp := entry.Time.Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf("%s %s %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)), nil

*/
func (formatter JSONLoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	time_formatter := "2006-01-02T15:04:05Z07:00.000"
	servicename := os.Getenv(envServiceName)
	field := map[string]interface{}{}
	data := map[string]interface{}{
		"servicename": servicename,
		"msg":         entry.Message,
		"level":       strings.ToUpper(entry.Level.String()),
		"time":        entry.Time.Format(time_formatter),
	}
	for key, val := range entry.Data {
		if xtools.Has(key, ctxParentID, ctxSpanID, ctxTraceID) {
			sval, isok := val.(string)
			if isok && len(sval) > 1 {
				data[key] = sval
			}
			continue
		}
		field[key] = val
	}

	if len(field) > 0 {
		data["entrydata"] = xtools.JSONToStr(field)
	}
	if entry.Caller != nil {
		data["file"] = fmt.Sprintf("%s:%s", entry.Caller.File, strconv.Itoa(entry.Caller.Line))
	}
	buf := bytes.NewBufferString(xtools.JSONMarshal(data))
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
