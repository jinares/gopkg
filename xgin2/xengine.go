package xgin2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	XEngine struct {
		*gin.Engine
	}
)

//NewXGin new enginx
func NewXGin() *XEngine {

	return &XEngine{gin.New()}
}

// ServeHTTP conforms to the http.Handler interface.
func (h *XEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	h.Engine.ServeHTTP(w, req)

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
