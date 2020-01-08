package xerr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jinares/gopkg/xtools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	XError struct {
		Code    XCode  `json:"code"`
		Message string `json:"message"`
		Where   string `json:"where"`
	}
)

func (se *XError) GRPCStatus() *status.Status {
	return status.New(codes.Code(se.Code), se.Message)
}
func caller(calldepth int, short bool) string {
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	} else if short {
		file = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d", file, line)
}

/*
Error() string
*/
func (p *XError) Error() string {
	return fmt.Sprintf("%d %s", (p.Code), p.Message)
}

/*
Error() string
*/
func (p *XError) String() string {
	return fmt.Sprintf("%s %s", p.Error(), p.Where)
}
func String(err error) string {
	if err == nil {
		return ""
	}
	if val, isok := (err).(*XError); isok {
		return val.String()
	}
	return err.Error()
}

func XErr(code XCode, msg string, where ...bool) *XError {

	if len(where) > 0 && where[0] {
		return &XError{
			Code:    code,
			Message: msg,
			Where:   caller(1, true),
		}
	}
	return &XError{
		Code:    code,
		Message: msg,
	}

}
func FromXErr(err error) *XError {
	if err == nil {
		return nil
	}
	if xe, isok := err.(*XError); isok {
		return xe
	}
	rpcerr, isok := status.FromError(err)
	if isok {
		return &XError{Code: XCode(rpcerr.Code()), Message: rpcerr.Message()}
	}
	msg := err.Error()
	index := strings.Index(msg, " ")
	if index < 1 {
		return &XError{
			Code:    Unknown,
			Message: msg,
		}
	}
	codestr := msg[0:index]

	icode, isok := xtools.IntVal(codestr)
	if isok == false {
		return &XError{
			Code:    Unknown,
			Message: msg,
		}
	}
	return &XError{Code: XCode(icode), Message: msg[index:]}
}
