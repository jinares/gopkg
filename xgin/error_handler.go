package xgin

import (
	"net/http"

	"github.com/jinares/gopkg/xlog"
	"github.com/jinares/gopkg/xtools"
)

func errorHandler() func(writer http.ResponseWriter, request *http.Request, err error) {

	return func(writer http.ResponseWriter, request *http.Request, err error) {
		if err != nil {
			xlog.CtxLog(request.Context()).Tracef("url:%s, errorhandler:%s", request.URL.Path, err.Error())
		}
		writer.WriteHeader(http.StatusOK)
		qu := xtools.UrlDecode(request.URL.RawQuery)
		callback := ""
		if qu, isok := qu["jsonp"]; isok && qu != "" {
			callback = qu
		}
		if qu, isok := qu["callback"]; isok && qu != "" {
			callback = qu
		}
		ret := map[string]interface{}{
			"msg": "服务繁忙,请稍后再试",
		}
		if callback != "" {
			writer.Write([]byte(callback + "(" + xtools.JSONToStr(ret) + ")"))
		} else {
			writer.Write([]byte(xtools.JSONToStr(ret)))
		}
	}
}
