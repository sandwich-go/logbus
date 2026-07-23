package logbus

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSetLogLevel(t *testing.T) {
	var buf bytes.Buffer

	Init(NewConf(
		WithLogLevel(zap.InfoLevel),
		WithWriteSyncer(zapcore.AddSync(&buf)),
		WithDisableTruncateWriteSyncer(true),
	))
	defer resetLogBus()
	defer Close()

	scopeLogger := NewScopeLogger("scope")

	Debug(testContext, "global debug before change")
	scopeLogger.Debug(testContext, "scope debug before change")

	SetLogLevel(zap.DebugLevel)

	Debug(testContext, "global debug after change")
	scopeLogger.Debug(testContext, "scope debug after change")

	output := buf.String()
	if strings.Contains(output, "before change") {
		t.Fatalf("debug logs should be filtered before SetLogLevel, output=%s", output)
	}
	if !strings.Contains(output, "global debug after change") {
		t.Fatalf("global debug log should be emitted after SetLogLevel, output=%s", output)
	}
	if !strings.Contains(output, "scope debug after change") {
		t.Fatalf("scope debug log should be emitted after SetLogLevel, output=%s", output)
	}
	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("runtime log level should be updated to debug, got %s", got.String())
	}
}

func TestSetLogLevelConcurrent(t *testing.T) {
	var buf bytes.Buffer

	Init(NewConf(
		WithLogLevel(zap.InfoLevel),
		WithWriteSyncer(zapcore.Lock(zapcore.AddSync(&buf))),
		WithDisableTruncateWriteSyncer(true),
	))
	defer resetLogBus()
	defer Close()

	scopeLogger := NewScopeLogger("scope")
	levels := []zapcore.Level{
		zap.DebugLevel,
		zap.InfoLevel,
		zap.WarnLevel,
		zap.ErrorLevel,
	}

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				level := levels[(worker+j)%len(levels)]
				SetLogLevel(level)
				Debug(testContext, "global concurrent debug")
				scopeLogger.Debug(testContext, "scope concurrent debug")
				Info(testContext, "global concurrent info")
				scopeLogger.Info(testContext, "scope concurrent info")
			}
		}(i)
	}

	close(start)
	wg.Wait()

	SetLogLevel(zap.DebugLevel)
	Debug(testContext, "global debug after concurrent change")
	scopeLogger.Debug(testContext, "scope debug after concurrent change")

	output := buf.String()
	if !strings.Contains(output, "global debug after concurrent change") {
		t.Fatalf("global debug log should be emitted after concurrent SetLogLevel, output=%s", output)
	}
	if !strings.Contains(output, "scope debug after concurrent change") {
		t.Fatalf("scope debug log should be emitted after concurrent SetLogLevel, output=%s", output)
	}
	if got := runtimeLogLevel.Level(); got != zap.DebugLevel {
		t.Fatalf("runtime log level should end at debug, got %s", got.String())
	}
	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("runtime log level should end at debug, got %s", got.String())
	}
}
