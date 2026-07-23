package logbus

import (
	"sync/atomic"

	"go.uber.org/zap/zapcore"
)

type configSnapshot struct {
	Conf
	writer  zapcore.WriteSyncer
	encoder zapcore.EncoderConfig
}

var activeConfig atomic.Pointer[configSnapshot]

func snapshotConf(conf *Conf) *configSnapshot {
	if conf == nil {
		panic("logbus: nil config")
	}

	cloned := cloneConf(conf)
	return &configSnapshot{Conf: cloned}
}

func cloneConf(conf *Conf) Conf {
	cloned := *conf
	cloned.DefaultPercentiles = append([]float64(nil), conf.DefaultPercentiles...)
	if conf.DefaultLabel != nil {
		cloned.DefaultLabel = make(map[string]string, len(conf.DefaultLabel))
		for key, value := range conf.DefaultLabel {
			cloned.DefaultLabel[key] = value
		}
	}
	if conf.TruncateWriteSyncerOption != nil {
		option := *conf.TruncateWriteSyncerOption
		cloned.TruncateWriteSyncerOption = &option
	}
	return cloned
}

func (c *configSnapshot) detached() *Conf {
	cloned := cloneConf(&c.Conf)
	return &cloned
}

func (c *configSnapshot) prepareWriter() {
	writer := c.WriteSyncer
	if c.BufferedStdout {
		BufferedWriteSyncer.WS = c.WriteSyncer
		writer = BufferedWriteSyncer
	}
	if !c.DisableTruncateWriteSyncer {
		if c.TruncateWriteSyncerOption == nil {
			panic("logbus: TruncateWriteSyncerOption cannot be nil when truncation is enabled")
		}
		writer = NewTruncateWriteSyncer(writer, c.TruncateWriteSyncerOption)
	}
	c.writer = writer
}

func currentConfig() *configSnapshot {
	if config := activeConfig.Load(); config != nil {
		return config
	}
	config := snapshotConf(newDefaultConf())
	config.encoder = EncodeConfig
	config.prepareWriter()
	if activeConfig.CompareAndSwap(nil, config) {
		return config
	}
	return activeConfig.Load()
}
