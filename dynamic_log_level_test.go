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

func writeOpsConfig(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s failed: %v", bizopsConfFile, err)
	}
}

func TestEnableDynamicLogLevel_EnvOnly(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{"env_config":{"log_level":"info"}}`)

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
	writeOpsConfig(t, confPath, `{"env_config":{"log_level":"warn"}}`)

	if !waitLogLevel(zap.WarnLevel, 3*time.Second) {
		t.Fatalf("expect level=warn after file change, got %s", GetLogLevel().String())
	}
}

func TestEnableDynamicLogLevel_ServiceOverridesEnv(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{
  "env_config": {"log_level": "info"},
  "service_config": {
    "svcA": {"log_level": "debug"}
  }
}`)

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
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{
  "env_config": {"log_level": "warn"},
  "service_config": {"svcA": {}}
}`)
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
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{"env_config":{"log_level":"not-a-level"}}`)

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

// TestEnableDynamicLogLevel_ServiceInvalidFallbackToEnv 验证 service 配置非法时回退到 env 级别。
// 场景：env_config.log_level=debug 合法；service_config.svcA.log_level=abc 非法。
// 期望：service_a 不被卡住，能用上 env 的 debug。
func TestEnableDynamicLogLevel_ServiceInvalidFallbackToEnv(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{
  "env_config": {"log_level": "debug"},
  "service_config": {"svcA": {"log_level": "not-a-level"}}
}`)
	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "svcA")

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("expect fallback to env debug when service invalid, got %s", got.String())
	}
}

// TestEnableDynamicLogLevel_AllInvalid 验证 service 与 env 均非法时保持原 level 不变。
func TestEnableDynamicLogLevel_AllInvalid(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{
  "env_config": {"log_level": "bad-env"},
  "service_config": {"svcA": {"log_level": "bad-svc"}}
}`)
	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "svcA")

	Init(NewConf(WithLogLevel(zap.InfoLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.InfoLevel {
		t.Fatalf("expect keep init info when all candidates invalid, got %s", got.String())
	}
}

// TestEnableDynamicLogLevel_EnvInvalidServiceValid 验证 env 非法 + service 合法时 service 生效。
func TestEnableDynamicLogLevel_EnvInvalidServiceValid(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{
  "env_config": {"log_level": "bogus"},
  "service_config": {"svcA": {"log_level": "debug"}}
}`)
	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "svcA")

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("expect service debug when env invalid, got %s", got.String())
	}
}

// TestCloseDynamicLogLevel_NoPanicOnNilDone 回归 xconf <= v0.3.28 缺陷：
// kv.Common.Done 在 New 时未 make，Close 中 close(c.Done) 直接 panic。
// xconf v0.3.29 已修复（New 中 make Done + Close 用 sync.Once 幂等）。
// 用例验证：在 v0.3.29 之上，单次 Close 与重复 Close 都不会 panic。
func TestCloseDynamicLogLevel_NoPanicOnNilDone(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{"env_config":{"log_level":"info"}}`)
	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "")

	Init(NewConf(WithLogLevel(zap.DebugLevel)))
	if !dynamicLogLevelStarted {
		t.Fatalf("dynamic loader should be started")
	}

	// 第一次 Close：原缺陷下会 panic（被 recover）。修复后应正常返回。
	closeDynamicLogLevel()
	if dynamicLogLevelStarted {
		t.Fatalf("dynamic loader should be stopped after Close")
	}

	// 第二次 Close：幂等性验证，不应 panic。
	closeDynamicLogLevel()
}

func TestEnableDynamicLogLevel_ReapplyAfterInit(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, bizopsConfFile)
	writeOpsConfig(t, confPath, `{"env_config":{"log_level":"debug"}}`)

	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "")

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("expect first Init apply dynamic debug, got %s", got.String())
	}

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	if got := GetLogLevel(); got != zap.DebugLevel {
		t.Fatalf("expect repeated Init reapply dynamic debug, got %s", got.String())
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
