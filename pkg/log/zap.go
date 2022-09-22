package log

import (
	"bj-pfd2/pkg/cfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var ZapLogger *zap.Logger             // 在每一微秒和每一次内存分配都很重要的上下文中，使用Logger，只支持强类型的结构化日志记录
var ZapSugarLogger *zap.SugaredLogger // 在性能很好但不是很关键的上下文中，使用SugaredLogger，支持结构化和 printf 风格
var zapLogLevel zapcore.LevelEnabler

func InitZapLogger(fName string, level string, maxSize int, maxNum int, maxDay int, compress bool) {
	lCfg = &logCfg{
		fileName:   fName,
		level:      level,
		maxSizeMb:  maxSize,
		maxFileNum: maxNum,
		maxFileDay: maxDay,
		compress:   compress,
	}

	switch lCfg.level {
	case "debug":
		zapLogLevel = zap.DebugLevel
	case "info":
		zapLogLevel = zap.InfoLevel
	case "warn":
		zapLogLevel = zap.WarnLevel
	case "error":
		zapLogLevel = zap.ErrorLevel
	case "panic":
		zapLogLevel = zap.PanicLevel
	case "fatal":
		zapLogLevel = zap.FatalLevel
	default:
		zapLogLevel = zap.InfoLevel
	}
	core := zapcore.NewTee(getAllZapCores()...)

	// AddCaller - 调用函数信息记录到日志中的功能
	// AddCallerSkip - 向上跳 1 层，输出调用封装日志函数的位置
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	ZapSugarLogger = ZapLogger.Sugar()
	defer ZapLogger.Sync() // flushes buffer, if any
	defer ZapSugarLogger.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 格式化时间显示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 使用大写字母记录日志级别
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := getLumberJackLogger(lCfg)
	return zapcore.AddSync(lumberJackLogger)
}

func getAllZapCores() []zapcore.Core {
	var allCore []zapcore.Core

	encoder := getEncoder()
	// debug 模式、显式输出到stdout 或 仅输出到 stdout 时将日志同时输出到 stdout
	if cfg.GetString("log.level") == "debug" ||
		cfg.GetBool("log.stdout") ||
		cfg.GetBool("log.only_stdout") {
		consoleDebugging := zapcore.Lock(os.Stdout)
		allCore = append(allCore, zapcore.NewCore(encoder, consoleDebugging, zapLogLevel))
	}
	// 仅输出到 stdout 时屏蔽文件输入
	if !cfg.GetBool("log.only_stdout") {
		writeSyncer := getZapLogWriter()
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapLogLevel))
	}
	return allCore
}
