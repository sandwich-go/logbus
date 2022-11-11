package basics

import (
	"io"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewFileWriter(tag string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filepath.Join(Setting.LocalLogDir, tag+".log"),
		MaxSize:    Setting.LocalLogMaxSize, // megabytes
		MaxBackups: Setting.LocalLogMaxBackups,
		MaxAge:     Setting.LocalLogMaxAge, // days
	}
}
