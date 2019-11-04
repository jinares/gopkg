package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

//InterceptorOpenTracing opentracing
func InterceptorOpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		tr := opentracing.GlobalTracer()
		ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)

		sp := tr.StartSpan(c.Request.URL.Path, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, c.Request.URL.Path)
		ext.Component.Set(sp, "net/http")
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), sp),
		)
		c.Next()
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
		sp.Finish()
	}
}
