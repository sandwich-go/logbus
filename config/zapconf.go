package config

import (
	"go.uber.org/zap"
)

var ZapConf = zap.Config{
	Development:      false,
	Encoding:         "json",
	EncoderConfig:    EncodeConfig,
	OutputPaths:      []string{},
	ErrorOutputPaths: []string{},
}
