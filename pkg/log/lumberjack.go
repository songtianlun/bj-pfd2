package log

import "github.com/natefinch/lumberjack"

func getLumberJackLogger(cfg *logCfg) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.fileName,
		MaxSize:    cfg.maxSizeMb,
		MaxBackups: cfg.maxFileNum,
		MaxAge:     cfg.maxFileDay,
		Compress:   cfg.compress,
	}
}
