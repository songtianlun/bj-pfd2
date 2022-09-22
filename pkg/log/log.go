package log

type logCfg struct {
	fileName   string
	level      string
	maxSizeMb  int
	maxFileNum int
	maxFileDay int
	compress   bool
}

var lCfg *logCfg

func InitLogger(fName string, level string, maxSize int, maxNum int, maxDay int, compress bool) {
	lCfg = &logCfg{
		fileName:   fName,
		level:      level,
		maxSizeMb:  maxSize,
		maxFileNum: maxNum,
		maxFileDay: maxDay,
		compress:   compress,
	}
	InitZapLogger(fName, level, maxSize, maxNum, maxDay, compress)
	InitLogRus()
}

func DebugF(format string, v ...interface{}) { ZapSugarLogger.Debugf(format, v...) }
func InfoF(format string, v ...interface{})  { ZapSugarLogger.Infof(format, v...) }
func WarnF(format string, v ...interface{})  { ZapSugarLogger.Warnf(format, v...) }
func ErrorF(format string, v ...interface{}) { ZapSugarLogger.Errorf(format, v...) }
func PanicF(format string, v ...interface{}) { ZapSugarLogger.Panicf(format, v...) }
func FatalF(format string, v ...interface{}) { ZapSugarLogger.Fatalf(format, v...) }
func Debug(v ...interface{})                 { ZapSugarLogger.Debug(v...) }
func Info(v ...interface{})                  { ZapSugarLogger.Info(v...) }
func Warn(v ...interface{})                  { ZapSugarLogger.Warn(v...) }
func Error(v ...interface{})                 { ZapSugarLogger.Error(v...) }
func Panic(v ...interface{})                 { ZapSugarLogger.Panic(v...) }
func Fatal(v ...interface{})                 { ZapSugarLogger.Fatal(v...) }
