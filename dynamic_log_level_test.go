package logbus

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// resetDynamicLogLevelForTest 清理 dynamic loader 状态，便于用例间复用。
func resetDynamicLogLevelForTest(t *testing.T) {
	t.Helper()
	closeDynamicLogLevel()
}

func writeBizopsYAML(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write bizops.yaml failed: %v", err)
	}
}

func TestEnableDynamicLogLevel_EnvOnly(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, logConfFile)
	writeBizopsYAML(t, confPath, "env_config:\n  log_level: info\n")

	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "")

	// Init 从默认 DebugLevel 开始
	Init(NewConf(WithLogLevel(zap.DebugLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	// 初始读取应已应用为 info
	if got := GetLogLevel(); got != zap.InfoLevel {
		t.Fatalf("expect initial level=info, got %s", got.String())
	}

	// 变更文件 -> warn
	writeBizopsYAML(t, confPath, "env_config:\n  log_level: warn\n")

	if !waitLogLevel(zap.WarnLevel, 3*time.Second) {
		t.Fatalf("expect level=warn after file change, got %s", GetLogLevel().String())
	}
}

func TestEnableDynamicLogLevel_ServiceOverridesEnv(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, logConfFile)
	writeBizopsYAML(t, confPath, `env_config:
  log_level: info
service_config:
  svcA:
    log_level: debug
`)

	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "svcA")

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	// svcA 覆盖为 debug
	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("expect initial level=debug (service override), got %s", got.String())
	}

	// 其他服务名时走 env 级 info
	closeDynamicLogLevel()
	t.Setenv(envKeyCDService, "svcB")
	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	if got := GetLogLevel(); got != zap.InfoLevel {
		t.Fatalf("expect level=info for svcB fallback, got %s", got.String())
	}
}

func TestEnableDynamicLogLevel_ServiceFieldLevelOverride(t *testing.T) {
	// 服务节点存在但未设置 log_level，应回落到 env 级。
	dir := t.TempDir()
	confPath := filepath.Join(dir, logConfFile)
	writeBizopsYAML(t, confPath, `env_config:
  log_level: warn
service_config:
  svcA: {}
`)
	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "svcA")

	Init(NewConf(WithLogLevel(zap.DebugLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.WarnLevel {
		t.Fatalf("expect fallback to env warn, got %s", got.String())
	}
}

func TestEnableDynamicLogLevel_NoEnv(t *testing.T) {
	// 未设置 sys_conf_path_env 时应静默跳过，初始 LogLevel 保持不变。
	t.Setenv(envKeyConfPathEnv, "")
	t.Setenv(envKeyCDService, "")

	Init(NewConf(WithLogLevel(zap.ErrorLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.ErrorLevel {
		t.Fatalf("expect level keep error when env not set, got %s", got.String())
	}
	if dynamicLogLevelStarted {
		t.Fatalf("dynamic loader should not start when env not set")
	}
}

func TestEnableDynamicLogLevel_InvalidLevel(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, logConfFile)
	writeBizopsYAML(t, confPath, "env_config:\n  log_level: not-a-level\n")

	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "")

	Init(NewConf(WithLogLevel(zap.InfoLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	// 非法 level 应被忽略，维持初始 info。
	if got := GetLogLevel(); got != zap.InfoLevel {
		t.Fatalf("expect level keep info on invalid input, got %s", got.String())
	}
}

func waitLogLevel(expect zapcore.Level, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if GetLogLevel() == expect {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return GetLogLevel() == expect
}
