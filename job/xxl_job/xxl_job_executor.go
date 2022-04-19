package xxl_job

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/william094/iris-lib"
	"go.uber.org/zap"
)

func InitExecutor(app *iris.Application, conf *iris_lib.Application) *Executor {
	//初始化执行器
	exec := NewExecutor(
		ServerAddr(conf.XXLJob.Addr),
		AccessToken(""),                                   //请求令牌(默认为空)
		ExecutorPort(fmt.Sprintf("%d", conf.Server.Port)), //默认9999（此处要与gin服务启动port必需一至）
		RegistryKey(conf.XXLJob.ExecutorName),             //执行器名称
		SetLogger(iris_lib.SystemLogger),
	)
	exec.Init()
	defer exec.Stop()
	//设置日志查看handler
	exec.LogHandler(func(req *LogReq) *LogRes {
		return &LogRes{Code: 200, Msg: "", Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})
	app.Post("run", func(ctx iris.Context) {
		exec.RunTask(ctx.ResponseWriter(), ctx.Request())
	})
	app.Post("kill", func(ctx iris.Context) {
		exec.KillTask(ctx.ResponseWriter(), ctx.Request())
	})
	app.Post("log", func(ctx iris.Context) {
		exec.TaskLog(ctx.ResponseWriter(), ctx.Request())
	})
	app.Post("beat", func(ctx iris.Context) {
		exec.TaskBeat(ctx.ResponseWriter(), ctx.Request())
	})
	app.Post("idleBeat", func(ctx iris.Context) {
		exec.TaskIdleBeat(ctx.ResponseWriter(), ctx.Request())
	})
	return &exec
}

type XxlJobLog struct {
}

func (l *XxlJobLog) Info(format string, a ...interface{}) {
	//logx.SystemLogger.Info(format, zap.Any("params", a))
}

func (l *XxlJobLog) Error(format string, a ...interface{}) {
	iris_lib.SystemLogger.Error(format, zap.Any("params", a))
}
