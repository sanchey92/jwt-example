package logger

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// TODO: add config (env) for logger!!!
var (
	logLevel     = "info"
	loggerFormat = "json"
	loggerOutput = []string{"stdout", "file"}

	fileName   = "logs/app.log"
	maxSize    = 10
	maxBackups = 3
	maxAge     = 7
)

func Init() {
	once.Do(func() {
		level := zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(logLevel)); err != nil {
			level.SetLevel(zap.InfoLevel)
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:      "timestamp",
			LevelKey:     "level",
			MessageKey:   "msg",
			CallerKey:    "caller",
			EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		}

		var encoder zapcore.Encoder

		if loggerFormat == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		var cores []zapcore.Core

		for _, output := range loggerOutput {
			switch strings.ToLower(output) {
			case "stdout":
				cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
			case "file":
				writer := zapcore.AddSync(&lumberjack.Logger{
					Filename:   fileName,
					MaxSize:    maxSize,
					MaxBackups: maxBackups,
					MaxAge:     maxAge,
					Compress:   true,
				})
				cores = append(cores, zapcore.NewCore(encoder, writer, level))
			}
		}

		logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
}

func GetLogger() *zap.Logger {
	if logger == nil {
		Init()
	}
	return logger
}
