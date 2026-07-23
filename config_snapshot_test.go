package logbus

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/sandwich-go/boost/xerror"
	"go.uber.org/zap/zapcore"
)

func isolateConfigTest(t *testing.T) {
	t.Helper()
	closeDynamicLogLevel()
	t.Setenv(envKeyConfPathEnv, "")
}

func TestLazyDefaultConfigCanBuildLogger(t *testing.T) {
	isolateConfigTest(t)
	previousConfig := activeConfig.Load()
	previousRefresh := refresh
	previousCache := cacheGLogger
	activeConfig.Store(nil)
	refresh = true
	cacheGLogger = nil
	t.Cleanup(func() {
		activeConfig.Store(previousConfig)
		refresh = previousRefresh
		cacheGLogger = previousCache
	})

	config := currentConfig()
	var output bytes.Buffer
	config.writer = zapcore.AddSync(&output)

	logger := NewScopeLogger("early")
	logger.Info(context.Background(), "before init")
	if err := logger.Sync(); err != nil {
		t.Fatalf("logger.Sync() error = %v", err)
	}

	got := output.String()
	for _, expected := range []string{
		`"` + LevelKey + `":"info"`,
		`"` + MsgBody + `":"before init"`,
		`"` + Tags + `":"early"`,
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("lazy default logger output does not contain %q: %s", expected, got)
		}
	}
}

func TestLoggerCapturesConfigAtConstruction(t *testing.T) {
	isolateConfigTest(t)
	var output bytes.Buffer
	conf := NewConf(
		WithDefaultChannel("original-channel"),
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "original")}
		}),
		WithWriteSyncer(zapcore.AddSync(&output)),
		WithDisableTruncateWriteSyncer(true),
	)

	Init(conf)
	t.Cleanup(func() {
		Close()
		Init(NewConf())
	})

	conf.DefaultChannel = "input-mutated-channel"
	conf.FetchLogContext = func() []Field {
		return []Field{String("config", "input-mutated")}
	}
	Setting.DefaultChannel = "setting-mutated-channel"
	Setting.FetchLogContext = func() []Field {
		return []Field{String("config", "setting-mutated")}
	}

	logger := NewScopeLogger("")
	SetGlobalGLogger(nil, "", false, 0)
	Info(context.Background(), "global")
	logger.Info(context.Background(), "scope")
	if err := logger.Sync(); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "original-channel") {
		t.Fatalf("output does not use snapshotted channel: %s", got)
	}
	if count := strings.Count(got, `"config":"original"`); count != 2 {
		t.Fatalf("snapshotted fetch field count = %d, want 2; output=%s", count, got)
	}
	for _, unexpected := range []string{"input-mutated", "setting-mutated"} {
		if strings.Contains(got, unexpected) {
			t.Fatalf("output contains post-Init config %q: %s", unexpected, got)
		}
	}
}

func TestScopeLoggerKeepsConfigAcrossInit(t *testing.T) {
	isolateConfigTest(t)
	var oldOutput bytes.Buffer
	Init(NewConf(
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "old")}
		}),
		WithWriteSyncer(zapcore.AddSync(&oldOutput)),
		WithDisableTruncateWriteSyncer(true),
	))
	oldLogger := NewScopeLogger("old")

	var newOutput bytes.Buffer
	Init(NewConf(
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "new")}
		}),
		WithWriteSyncer(zapcore.AddSync(&newOutput)),
		WithDisableTruncateWriteSyncer(true),
	))
	t.Cleanup(func() {
		Close()
		Init(NewConf())
	})
	newLogger := NewScopeLogger("new")

	oldLogger.Info(context.Background(), "old logger")
	oldLogger.DebugWithChannel(context.Background(), "old-channel", "old channel logger")
	newLogger.Info(context.Background(), "new logger")
	newLogger.DebugWithChannel(context.Background(), "new-channel", "new channel logger")
	if err := oldLogger.Sync(); err != nil {
		t.Fatalf("oldLogger.Sync() error = %v", err)
	}
	if err := newLogger.Sync(); err != nil {
		t.Fatalf("newLogger.Sync() error = %v", err)
	}

	if got := oldOutput.String(); !strings.Contains(got, `"config":"old"`) || strings.Contains(got, `"config":"new"`) {
		t.Fatalf("old logger did not keep old config: %s", got)
	} else if !strings.Contains(got, "old-channel") || strings.Contains(got, "new-channel") {
		t.Fatalf("old logger channel output used a different logger: %s", got)
	}
	if got := newOutput.String(); !strings.Contains(got, `"config":"new"`) || strings.Contains(got, `"config":"old"`) {
		t.Fatalf("new logger did not use new config: %s", got)
	} else if !strings.Contains(got, "new-channel") || strings.Contains(got, "old-channel") {
		t.Fatalf("new logger channel output used a different logger: %s", got)
	}
}

func TestLogicalErrorPolicyUsesInitConfig(t *testing.T) {
	isolateConfigTest(t)
	logical := xerror.NewText("logical").SetLogic()

	Init(NewConf(
		WithIgnoreLogicalError(true),
		WithDisableTruncateWriteSyncer(true),
	))
	t.Cleanup(func() {
		Close()
		Init(NewConf())
	})

	if field := ErrorField(logical); field.Type != zapcore.StringType {
		t.Fatalf("ErrorField(logical).Type = %v, want string", field.Type)
	}

	Setting.IgnoreLogicalError = false
	if field := ErrorField(logical); field.Type != zapcore.StringType {
		t.Fatalf("Setting mutation changed logical error policy: %v", field.Type)
	}
	Init(NewConf(
		WithIgnoreLogicalError(false),
		WithDisableTruncateWriteSyncer(true),
	))
	if field := ErrorField(logical); field.Type != zapcore.ErrorType {
		t.Fatalf("ErrorField(logical).Type = %v, want error", field.Type)
	}
}

func TestEarlyScopeLoggerRefreshesConfigSnapshot(t *testing.T) {
	isolateConfigTest(t)
	var oldOutput bytes.Buffer
	Init(NewConf(
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "old")}
		}),
		WithWriteSyncer(zapcore.AddSync(&oldOutput)),
		WithDisableTruncateWriteSyncer(true),
	))

	refresh = true
	cacheGLogger = nil
	defaultLogger := NewScopeLogger("default")
	nilOverride := NewScopeLoggerWithFetchFunc("nil", nil)
	customLogger := NewScopeLoggerWithFetchFunc("custom", func() []Field {
		return []Field{String("config", "custom")}
	})
	inheritLogger := NewScopeLogger("inherit")
	forcedLogger := NewScopeLoggerPrintAsError("forced")
	defaultLogger.GInfoDepth(context.Background(), 0, "populate old depth cache")

	var newOutput bytes.Buffer
	Init(NewConf(
		WithPrintAsError(false),
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "new")}
		}),
		WithWriteSyncer(zapcore.AddSync(&newOutput)),
		WithDisableTruncateWriteSyncer(true),
	))
	t.Cleanup(func() {
		Close()
		Init(NewConf())
	})

	defaultLogger.GInfoDepth(context.Background(), 0, "refreshed depth logger")
	nilOverride.Info(context.Background(), "nil override")
	customLogger.Info(context.Background(), "custom override")
	inheritLogger.Debug(context.Background(), "inherit print mode", ErrorField(errors.New("inherit")))
	forcedLogger.Debug(context.Background(), "forced print mode", ErrorField(errors.New("forced")))
	if err := defaultLogger.Sync(); err != nil {
		t.Fatalf("defaultLogger.Sync() error = %v", err)
	}

	got := newOutput.String()
	if count := strings.Count(got, `"config":"new"`); count != 4 {
		t.Fatalf("refreshed default fetch count = %d, want 4; output=%s", count, got)
	}
	if count := strings.Count(got, `"config":"custom"`); count != 1 {
		t.Fatalf("preserved custom fetch count = %d, want 1; output=%s", count, got)
	}
	if strings.Contains(got, `"config":"old"`) {
		t.Fatalf("refreshed logger still uses old config: %s", got)
	}
	assertLogLevel(t, got, "inherit print mode", "debug")
	assertLogLevel(t, got, "forced print mode", "error")
}

func assertLogLevel(t *testing.T, output, message, level string) {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, `"msg":"`+message+`"`) {
			if !strings.Contains(line, `"log_level":"`+level+`"`) {
				t.Fatalf("log %q does not have level %q: %s", message, level, line)
			}
			return
		}
	}
	t.Fatalf("log %q not found in output: %s", message, output)
}

func TestConfigSnapshotClonesMutableValues(t *testing.T) {
	percentiles := []float64{0.5, 0.9}
	labels := map[string]string{"source": "original"}
	truncate := NewTruncateWriteSyncerOption(WithTruncateMaxSize(100))
	conf := NewConf(
		WithDefaultPercentiles(percentiles...),
		WithDefaultLabel(labels),
		WithTruncateWriteSyncerOption(truncate),
	)

	snapshot := snapshotConf(conf)
	percentiles[0] = 1
	labels["source"] = "mutated"
	truncate.TruncateMaxSize = 200

	if got := snapshot.DefaultPercentiles[0]; got != 0.5 {
		t.Fatalf("snapshot percentile = %v, want 0.5", got)
	}
	if got := snapshot.DefaultLabel["source"]; got != "original" {
		t.Fatalf("snapshot label = %q, want original", got)
	}
	if got := snapshot.TruncateWriteSyncerOption.TruncateMaxSize; got != 100 {
		t.Fatalf("snapshot truncate max size = %d, want 100", got)
	}
}

func TestLoggerSnapshotDoesNotReadInputConfig(t *testing.T) {
	isolateConfigTest(t)
	var output bytes.Buffer
	conf := NewConf(
		WithFetchLogContext(func() []Field {
			return []Field{String("config", "snapshot")}
		}),
		WithWriteSyncer(zapcore.Lock(zapcore.AddSync(&output))),
		WithDisableTruncateWriteSyncer(true),
	)
	Init(conf)
	t.Cleanup(func() {
		Close()
		Init(NewConf())
	})
	logger := NewScopeLogger("race")

	const (
		writers = 4
		logs    = 100
	)
	start := make(chan struct{})
	var wait sync.WaitGroup
	wait.Add(writers + 1)
	go func() {
		defer wait.Done()
		<-start
		for i := 0; i < writers*logs*2; i++ {
			conf.FetchLogContext = func() []Field {
				return []Field{String("config", "mutated")}
			}
		}
	}()
	for worker := 0; worker < writers; worker++ {
		go func() {
			defer wait.Done()
			<-start
			for i := 0; i < logs; i++ {
				logger.Info(context.Background(), "concurrent")
			}
		}()
	}
	close(start)
	wait.Wait()

	got := output.String()
	if count := strings.Count(got, `"config":"snapshot"`); count != writers*logs {
		t.Fatalf("snapshot field count = %d, want %d", count, writers*logs)
	}
	if strings.Contains(got, `"config":"mutated"`) {
		t.Fatalf("logger read mutated input config")
	}
}
