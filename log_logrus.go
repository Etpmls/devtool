package d

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Log log

const (
	defaultPath = "log/app.log"
)

type log struct {
	Optional optionalLogrus
}

type optionalLogrus struct {
	Level      logrus.Level
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// 默认配置
// https://github.com/sirupsen/logrus
func (this *log) Init() {
	if this.Optional.Level == 0 {
		this.Optional.Level = logrus.InfoLevel
	}
	if this.Optional.Filename == "" {
		this.Optional.Filename = defaultPath
	}
	if this.Optional.MaxSize == 0 {
		this.Optional.MaxSize = 500
	}
	if this.Optional.MaxBackups == 0 {
		this.Optional.MaxBackups = 3
	}
	if this.Optional.MaxAge == 0 {
		this.Optional.MaxAge = 30
	}

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(this.GetRollingLogger())

	// Only logrus the warning severity or above.
	logrus.SetLevel(this.Optional.Level)

	// Logging Method Name
	logrus.SetReportCaller(true)

	// Hook
	// https://github.com/sirupsen/logrus/tree/master/hooks/writer
	logrus.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	logrus.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
}

// 日志分割
// https://github.com/natefinch/lumberjack
func (this *log) GetRollingLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   this.Optional.Filename,
		MaxSize:    this.Optional.MaxSize, // megabytes
		MaxBackups: this.Optional.MaxBackups,
		MaxAge:     this.Optional.MaxAge, //days
		Compress:   this.Optional.Compress, // disabled by default
	}
}






