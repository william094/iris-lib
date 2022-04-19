package iris_lib

import (
	"context"
	"github.com/hashicorp/go-uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

var (
	SystemLogger *zap.Logger
)

func InitLog(config *Application) *zap.Logger {
	logPath := config.Logger.FilePath
	applicationName := config.Logger.FileName
	enableConsoleOutPut := config.Logger.ConsoleEnable
	SystemLogger = globalLogZap(logPath, applicationName+"-"+"system", enableConsoleOutPut)
	return globalLogZap(logPath, applicationName, enableConsoleOutPut)
}

func WithContext(ctx context.Context) *zap.Logger {
	requestId := ctx.Value("TrackId")
	uid := ctx.Value("uid")
	logger := ctx.Value("log").(*zap.Logger)
	if requestId == nil {
		requestId, _ = uuid.GenerateUUID()
	}
	logger = logger.WithOptions(zap.Fields(zap.String("TrackId", requestId.(string))))
	if uid != nil {
		logger = logger.WithOptions(zap.Fields(zap.Int64("uid", uid.(int64))))
	}
	return logger
}

func globalLogZap(logPath, applicationName string, enableConsoleOutPut bool) *zap.Logger {
	infoPath := logPath + "/" + applicationName + "-" + "info.log"
	errorPath := logPath + "/" + applicationName + "-" + "error.log"
	encode := newEncoder()
	// 设置日志级别
	// 实现两个判断日志等级的interface (其实 zapcore.Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encode, zapcore.AddSync(NewWriter(infoPath)), infoLevel),
		zapcore.NewCore(encode, zapcore.AddSync(NewWriter(errorPath)), warnLevel),
	)
	if enableConsoleOutPut {
		core = zapcore.NewTee(
			zapcore.NewCore(encode, zapcore.AddSync(io.Writer(os.Stdout)), debugLevel),
			core,
		)
	}
	// 构造日志 开启开发模式，堆栈跟踪 开启文件及行号 跳过当前调用方，防止将封装方法当作调用方
	return zap.New(core, zap.AddCaller(), zap.Development(), zap.AddCallerSkip(1))
}

func newEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logx",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, // 时间格式化
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器，用于显示Method
		EncodeName:     zapcore.FullNameEncoder,
	})
}

func NewWriter(logName string) io.Writer {
	writer, err := rotatelogs.New(
		// 日志文件
		strings.Replace(logName, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(logName),
		// 日志周期(默认每86400秒/一天旋转一次)
		rotatelogs.WithRotationTime(time.Hour*24),
		// 清除历史 (WithMaxAge和WithRotationCount只能选其一)
		rotatelogs.WithMaxAge(time.Hour*24*30),     //每30天清除下日志文件
		rotatelogs.WithRotationSize(512*1024*1024), //日志文件大小，单位byte
		//rotatelogs.WithRotationCount(10), //只保留最近的N个日志文件
	)
	if err != nil {
		panic(err)
	}
	return writer
}
