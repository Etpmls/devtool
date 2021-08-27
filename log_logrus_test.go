package d_test

import (
	d "github.com/Etpmls/devtool"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	d.Log.Optional.Compress = true
	d.Log.Init()

	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Info("Panic Recover")
		}
	}()
	log.Panic("Panic")
}