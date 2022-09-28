package log

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"os"
)

type logrusAdapt struct {
	l *logrus.Logger
}

func (l *logrusAdapt) Debugf(format string, v ...interface{}) { l.l.Debugf(format, v...) }
func (l *logrusAdapt) Infof(format string, v ...interface{})  { l.l.Infof(format, v...) }
func (l *logrusAdapt) Warnf(format string, v ...interface{})  { l.l.Warnf(format, v...) }
func (l *logrusAdapt) Errorf(format string, v ...interface{}) { l.l.Errorf(format, v...) }
func (l *logrusAdapt) Panicf(format string, v ...interface{}) { l.l.Panicf(format, v...) }
func (l *logrusAdapt) Fatalf(format string, v ...interface{}) { l.l.Fatalf(format, v...) }
func (l *logrusAdapt) Debug(v ...interface{})                 { l.l.Debug(v...) }
func (l *logrusAdapt) Info(v ...interface{})                  { l.l.Info(v...) }
func (l *logrusAdapt) Warn(v ...interface{})                  { l.l.Warn(v...) }
func (l *logrusAdapt) Error(v ...interface{})                 { l.l.Error(v...) }
func (l *logrusAdapt) Panic(v ...interface{})                 { l.l.Panic(v...) }
func (l *logrusAdapt) Fatal(v ...interface{})                 { l.l.Fatal(v...) }

func newLogrus(lCfg *CfgLog) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	//l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.DebugLevel)

	l.AddHook(&writer.Hook{
		Writer:    getLumberJackLogger(lCfg),
		LogLevels: logrus.AllLevels,
	})

	l.Error("error")
	l.Warn("warn")
	l.Info("info")
	l.Debug("debug")

	return l
}

func InitLogrus(lCfg *CfgLog) Logger {
	return &logrusAdapt{
		l: newLogrus(lCfg),
	}
}
