package xgin

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinares/gopkg/xlog"
)

type (
	XEngine struct {
		*gin.Engine
		location *Locations
	}
)

//NewXGin new enginx
func NewXGin() *XEngine {

	return &XEngine{gin.New(), nil}
}

//SetLoction  set loction field
func (h *XEngine) SetLoction(locations *Locations) {
	if locations == nil {
		return
	}
	h.location = locations
}
func (h *XEngine) proxy(w http.ResponseWriter, r *http.Request) bool {
	if h.location == nil {
		return false
	}
	rule := h.location.Matching(r)
	if rule == nil {
		return false
	}

	if rule.IsProxy() == false {
		if rule.Rewrite != "" {
			r.URL.Path = strings.Replace(r.URL.Path, rule.GetLocationVal(), rule.Rewrite, 1)
		}
		h.Engine.ServeHTTP(w, r)
		return true
	}
	remote, err := url.Parse(rule.ProxyPass)
	if err != nil {
		xlog.CtxLog(r.Context()).WithField("conf", rule.ProxyPass).Error(err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	if rule.ProxyHost != "" {
		r.Host = rule.ProxyHost

	}
	proxy.ErrorHandler = errorHandler()
	if rule.ProxyRewrite != "" {
		r.URL.Path = strings.Replace(r.URL.Path, rule.GetLocationVal(), rule.ProxyRewrite, 1)
	}
	proxy.ServeHTTP(w, r)
	return true
}

// ServeHTTP conforms to the http.Handler interface.
func (h *XEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//fmt.Println(req.URL.Path)
	ti := time.Now()
	if h.proxy(w, req) {
		return
	}
	h.Engine.ServeHTTP(w, req)
	elapsed := time.Since(ti)
	if elapsed > (1 * time.Second) {
		xlog.CtxLog(req.Context()).WithField("slow", elapsed).Infof("url:%s method:%s", req.RequestURI, req.Method)
	}
}

//Run start
func (h *XEngine) Run(addr string) error {
	fmt.Println(fmt.Sprintf("Listening and serving HTTP on %s", addr))
	return http.ListenAndServe(addr, h)
}

//RunTLS start
func (h *XEngine) RunTLS(addr, certFile, keyFile string) (err error) {

	err = http.ListenAndServeTLS(addr, certFile, keyFile, h)
	return
}

//"test,omitempty
// regexp.MatchString(h.User, val.User)"
