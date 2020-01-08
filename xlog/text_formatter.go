package xlog

import (
	"bytes"
	"github.com/jinares/gopkg/xtools"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type TextLoggerFormatter struct {
}

/*

FormatDefault = "[${serviceName},${traceID},${spanID},${parentSpanID}][${pid}]"

timestamp := entry.Time.Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf("%s %s %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)), nil

*/
func (formatter TextLoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	time_formatter := "2006-01-02T15:04:05Z07:00.000"
	buffer := bytes.NewBufferString(entry.Time.Format(time_formatter))
	buffer.WriteString(" ")
	buffer.WriteString(strings.ToUpper(entry.Level.String()))
	buffer.WriteString(" ")
	buffer.WriteByte('[')
	buffer.WriteString(os.Getenv(envServiceName))
	buffer.WriteByte(',')
	for _, val := range []string{ctxTraceID, ctxSpanID, ctxParentID} {
		vval := entry.Data[val]
		sval, isok := vval.(string)
		if isok {
			buffer.WriteString(sval)
		}
		buffer.WriteByte(',')
	}

	buffer.WriteByte(']')
	buffer.WriteByte('[')
	buffer.WriteString(strconv.Itoa(os.Getpid()))
	buffer.WriteByte(']')
	buffer.WriteByte(' ')
	buffer.WriteString(entry.Message)
	//buffer.WriteByte('\n')

	for key, val := range entry.Data {
		if xtools.Has(key, ctxParentID, ctxSpanID, ctxTraceID) {
			continue
		}
		buffer.WriteByte('\n')
		buffer.WriteString(key)
		buffer.WriteString(" =  ")
		switch vv := val.(type) {
		case string:
			buffer.WriteString(vv)
		default:
			buffer.WriteString(xtools.JSONToStr(vv))
		}

	}

	if entry.Caller != nil {
		buffer.WriteByte('\n')
		buffer.WriteString(entry.Caller.File)
		buffer.WriteByte(':')
		buffer.WriteString(strconv.Itoa(entry.Caller.Line))
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}
