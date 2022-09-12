package plugin

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"net/http"
	"runtime"
)

func GlobalRecover(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			errMsg := fmt.Sprintf("错误信息: %s", err)
			// when stack finishes
			logMessage := fmt.Sprintf("从错误中回复：('%s')\n", ctx.HandlerName())
			logMessage += errMsg + "\n"
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			// 打印错误日志
			ctx.Application().Logger().Error(logMessage)
			// 返回错误信息
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
	}()
	ctx.Next()
}
