package d

import (
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Log logrus

const (
	defaultPath = "log/app.log"
)

type logrus struct {
	Level log.Level
	Filename string
	MaxSize int
	MaxBackups int
	MaxAge int
	Compress bool
}

// 默认配置
// https://github.com/sirupsen/logrus
func (this *logrus) Init() {
	if this.Level == 0 {
		this.Level = log.InfoLevel
	}
	if this.Filename == "" {
		this.Filename = defaultPath
	}
	if this.MaxSize == 0 {
		this.MaxSize = 500
	}
	if this.MaxBackups == 0 {
		this.MaxBackups = 3
	}
	if this.MaxAge == 0 {
		this.MaxAge = 30
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(this.GetRollingLogger())

	// Only log the warning severity or above.
	log.SetLevel(this.Level)

	// Logging Method Name
	log.SetReportCaller(true)

	// Hook
	// https://github.com/sirupsen/logrus/tree/master/hooks/writer
	log.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	})
}

// 日志分割
// https://github.com/natefinch/lumberjack
func (this *logrus) GetRollingLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   this.Filename,
		MaxSize:    this.MaxSize, // megabytes
		MaxBackups: this.MaxBackups,
		MaxAge:     this.MaxAge, //days
		Compress:   this.Compress, // disabled by default
	}
}






