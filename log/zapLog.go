package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)


var core zapcore.Core

var Logger *zap.Logger

var size uint16

func init() {
	// 设置一些基本日志格式
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	printLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	// 实现两个判断日志等级的interface
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.DebugLevel
	//})

	// 获取 info、error日志文件的io.Writer
	//infoWriter := getWriter("")

	// 创建Logger
	core = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), printLevel),   //打印到控制台
		//zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),   //普通日志
		//zapcore.NewCore(encoder, zapcore.AddSync(&LogWriter{}), infoLevel), //es日志
	)
	//分配log
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)) // 显示打日志点的文件名和行数
}

type LogWriter struct {
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	//AnalysisLog(string(p))
	return
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+"%Y%m%d"+".log",
		rotatelogs.WithLinkName(filename),                        //设置文件名称
		rotatelogs.WithMaxAge(time.Hour*24*7), //日志保存时间天
		//rotatelogs.WithRotationTime(time.Hour*24), //分割一次日志时间小时
		rotatelogs.WithRotationSize(1024*1024*100), //每个日志文件最大存储容量M
	)
	size++
	if err != nil {
		panic(err)
	}
	return hook
}

func GetCore() zapcore.Core {
	return core
}

func DefaultLogger() *zap.SugaredLogger {
	return Logger.Sugar()
}
