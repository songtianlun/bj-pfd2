package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 格式化时间显示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 使用大写字母记录日志级别
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLogWriter(lCfg *CfgLog) zapcore.WriteSyncer {
	lumberJackLogger := getLumberJackLogger(lCfg)
	return zapcore.AddSync(lumberJackLogger)
}

func getZapLogLevel(lv string) zapcore.LevelEnabler {
	switch lv {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func getAllZapCores(lCfg *CfgLog) []zapcore.Core {
	var allCore []zapcore.Core
	zapLogLevel := getZapLogLevel(lCfg.Level)

	encoder := getEncoder()
	// debug 模式、显式输出到stdout 或 仅输出到 Stdout 时将日志同时输出到 Stdout
	if lCfg.Level == "debug" ||
		lCfg.Stdout ||
		lCfg.OnlyStdout {
		consoleDebugging := zapcore.Lock(os.Stdout)
		allCore = append(allCore, zapcore.NewCore(encoder, consoleDebugging, zapLogLevel))
	}
	// 仅输出到 Stdout 时屏蔽文件输入
	if !lCfg.OnlyStdout {
		writeSyncer := getZapLogWriter(lCfg)
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapLogLevel))
	}
	return allCore
}

type zapAdapt struct {
	zl *zap.SugaredLogger // 在性能很好但不是很关键的上下文中，使用SugaredLogger，支持结构化和 printf 风格
}

func (z *zapAdapt) Debugf(format string, v ...interface{}) { z.zl.Debugf(format, v...) }
func (z *zapAdapt) Infof(format string, v ...interface{})  { z.zl.Infof(format, v...) }
func (z *zapAdapt) Warnf(format string, v ...interface{})  { z.zl.Warnf(format, v...) }
func (z *zapAdapt) Errorf(format string, v ...interface{}) { z.zl.Errorf(format, v...) }
func (z *zapAdapt) Panicf(format string, v ...interface{}) { z.zl.Panicf(format, v...) }
func (z *zapAdapt) Fatalf(format string, v ...interface{}) { z.zl.Fatalf(format, v...) }
func (z *zapAdapt) Debug(v ...interface{})                 { z.zl.Debug(v...) }
func (z *zapAdapt) Info(v ...interface{})                  { z.zl.Info(v...) }
func (z *zapAdapt) Warn(v ...interface{})                  { z.zl.Warn(v...) }
func (z *zapAdapt) Error(v ...interface{})                 { z.zl.Error(v...) }
func (z *zapAdapt) Panic(v ...interface{})                 { z.zl.Panic(v...) }
func (z *zapAdapt) Fatal(v ...interface{})                 { z.zl.Fatal(v...) }

func newZapAdapter(lCfg *CfgLog) (ZapSugarLogger *zap.SugaredLogger) {
	var ZapLogger *zap.Logger // 在每一微秒和每一次内存分配都很重要的上下文中，使用Logger，只支持强类型的结构化日志记录
	core := zapcore.NewTee(getAllZapCores(lCfg)...)

	// AddCaller - 调用函数信息记录到日志中的功能
	// AddCallerSkip - 向上跳 1 层，输出调用封装日志函数的位置
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	ZapSugarLogger = ZapLogger.Sugar()
	//defer func(z *zap.Logger) {
	//	err := z.Sync()
	//	if err != nil {
	//		//panic("zap logger sync error - " + err.Error())
	//	}
	//}(ZapLogger) // flushes buffer, if any
	//defer func(zs *zap.SugaredLogger) {
	//	err := zs.Sync()
	//	if err != nil {
	//		//panic("zap sugared logger sync error - " + err.Error())
	//	}
	//}(ZapSugarLogger)
	return
}

func InitZapAdapter(lCfg *CfgLog) Logger {
	return &zapAdapt{
		zl: newZapAdapter(lCfg),
	}
}
