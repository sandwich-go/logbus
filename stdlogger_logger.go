package logbus

// getLogger get StdLogger by tagKey. If StdLogger doesn't exist, add a new one
/*func getLogger(tags ...string) *StdLogger {
	tagKey := getTagKey(tags)

	if l, ok := loggerMap.Load(tagKey); ok {
		return l.(*StdLogger)
	}

	return addStdLogger(tagKey, tags...)
}

func getTagKey(tags []string) string {
	if len(tags) == 0 {
		tags = []string{DefaultTag}
	}
	for k, v := range tags {
		tags[k] = strings.TrimSpace(v)
	}
	sort.Strings(tags)

	return strings.Join(tags, ".")
}

// addStdLogger add a new StdLogger
func addStdLogger(tagKey string, tags ...string) *StdLogger {
	var cores []zapcore.Core
	encoder := newJSONEncoder(EncodeConfig)
	if Setting.Dev {
		encoder = newConsoleEncoder(EncodeConfig)
	}

	// stdout 只能输出到stdout
	var writer zapcore.WriteSyncer
	writer = os.Stdout
	if Setting.BufferedStdout {
		writer = BufferedWriteSyncer
	}
	stdCore := zapcore.NewCore(encoder, writer, Setting.LogLevel).With([]zap.Field{zap.String(Tags, tagKey)})
	cores = append(cores, stdCore)

	l := &StdLogger{
		// clone a new zapLogger with options
		z: gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		})),
		tags: tags,
	}

	loggerMap.Store(tagKey, l)

	return l
}*/
