package job

import (
	"github.com/kataras/iris/v12"
	"github.com/william094/iris-lib/configuration"
	"github.com/william094/iris-lib/logx"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"
)

func InitExecutor(app *iris.Application, conf *configuration.Application) xxl.Executor {
	//初始化执行器
	exec := xxl.NewExecutor(
		xxl.ServerAddr(conf.XXLJob.Addr),
		xxl.AccessToken(conf.XXLJob.AccessToken),
		xxl.ExecutorPort(conf.XXLJob.Port),
		xxl.RegistryKey(conf.XXLJob.ExecutorName),
		xxl.SetLogger(&XxlJobLog{}),
	)
	exec.Init()
	defer exec.Stop()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
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
		exec.Beat(ctx.ResponseWriter(), ctx.Request())
	})
	app.Post("idleBeat", func(ctx iris.Context) {
		exec.IdleBeat(ctx.ResponseWriter(), ctx.Request())
	})
	return exec
}

type XxlJobLog struct {
}

func (l *XxlJobLog) Info(format string, a ...interface{}) {
	//logx.SystemLogger.Info(format, zap.Any("params", a))
}

func (l *XxlJobLog) Error(format string, a ...interface{}) {
	logx.SystemLogger.Error(format, zap.Any("params", a))
}
