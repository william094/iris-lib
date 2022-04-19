package iris_lib

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

func StartServer(app *iris.Application, conf *Application) {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Server.Port),
		Handler:        app,
		ReadTimeout:    conf.Server.ReadTimeout * time.Second,
		WriteTimeout:   conf.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: conf.Server.MaxHeaderBytes,
	}
	fmt.Sprintf("Server Start  port %d env %s", conf.Server.Port, conf.Server.Environment)
	if err := app.Run(iris.Server(s), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations,
		iris.WithConfiguration(iris.Configuration{
			DisableInterruptHandler:           true,
			EnableOptimizations:               false,
			DisableBodyConsumptionOnUnmarshal: true,
			Charset:                           "UTF-8",
		})); err != nil {
		panic(err)
	}

}

func CloseServer(app *iris.Application) {
	fmt.Print("system shutdown signal")
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 关闭所有主机
	app.Shutdown(ctx)
	fmt.Print("http server shutdown")
}
