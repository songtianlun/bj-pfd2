package log

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"os"
)

func InitLogRus() {
	logrus.SetOutput(getLumberJackLogger(lCfg))
	//logrus.SetOutput(os.Stdout)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.AddHook(&writer.Hook{
		Writer:    os.Stdout,
		LogLevels: logrus.AllLevels,
	})

	//logrus.Fatal("Fatal")
	//logrus.Panic("panic")
	logrus.Error("error")
	logrus.Warn("warn")
	logrus.Info("info")
	logrus.Debug("debug")
}
