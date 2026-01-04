package pretty

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFileConfig struct {
	MaxSize    int  `default:"10"`
	MaxBackups int  `default:"5"`
	MaxAge     int  `default:"31"`
	Compress   bool `default:"true"`
}

func NewLumberjackLogger(logFileName string, config LogFileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

func DefaultLogFileConfig() LogFileConfig {
	return LogFileConfig{
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     31,
		Compress:   true,
	}
}

func NewLogFileConfig(maxSize, maxBackups, maxAge int, compress bool) LogFileConfig {
	return LogFileConfig{
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
