package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	if (os.Getenv("DEBUG") != "") && (os.Getenv("DEBUG") != "0") {
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg.Level.SetLevel(zap.InfoLevel)
	}

	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")
	log, _ := cfg.Build()
	l := log.Sugar()
	defer l.Sync()
	return l
}
