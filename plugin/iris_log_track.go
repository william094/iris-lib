package plugin

import (
	"bytes"
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/kataras/iris/v12"
	"github.com/william094/iris-lib/logx"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

func TracingHandler(ctx iris.Context) {
	var traceId string
	var ct context.Context
	traceId = ctx.Request().Header.Get("TrackId")
	if traceId == "" {
		if requestId := ctx.Request().Context().Value("TrackId"); requestId != nil {
			traceId = requestId.(string)
		} else {
			traceId, _ = uuid.GenerateUUID()
		}
	}
	ct = context.WithValue(ctx.Request().Context(), "TrackId", traceId)
	ctx.ResetRequest(ctx.Request().Clone(ct))
	ctx.Next()

}

// RequestLogPlugin 日志中间件
func RequestLogPlugin(ctx iris.Context) {
	path := ctx.Request().URL.RequestURI()
	method := ctx.Request().Method
	IP := ctx.GetHeader("X-Real-Ip")

	params := ""
	// 如果是POST/PUT请求，并且内容类型为JSON，则读取内容体
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err == nil {
			defer ctx.Request().Body.Close()
			buf := bytes.NewBuffer(body)
			ctx.Request().Body = ioutil.NopCloser(buf)
			params = string(body)
			if strings.Contains(params, "\r\n") {
				params = strings.ReplaceAll(params, "\r\n", "")
			}
			if strings.Contains(params, "\n") {
				params = strings.ReplaceAll(params, "\n", "")
			}
			params = strings.ReplaceAll(params, " ", "")
		}
	} else {
		params = ctx.Request().Form.Encode()
	}
	ctx.Next()
	//下面是返回日志
	respStr := string(ctx.Recorder().Body())
	logx.SystemLogger.Info("访问日志", zap.String("path", path),
		zap.String("method", method), zap.String("IP", IP), zap.String("params", params), zap.String("resp", respStr))
}
