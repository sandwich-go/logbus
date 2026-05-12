package logbus

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
)

// TestEnableDynamicLogLevel_FileCreatedLater 验证 Init 时文件不存在，之后被创建能否生效。
// 典型场景：K8s configmap 未绑定 / PMT 未下发 bizops.yaml，稍后下发。
func TestEnableDynamicLogLevel_FileCreatedLater(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, logConfFile)

	// 确保文件一开始不存在
	if _, err := os.Stat(confPath); !os.IsNotExist(err) {
		t.Fatalf("file should not exist initially: %v", err)
	}

	t.Setenv(envKeyConfPathEnv, dir)
	t.Setenv(envKeyCDService, "")

	Init(NewConf(WithLogLevel(zap.WarnLevel)))
	defer resetDynamicLogLevelForTest(t)
	defer Close()

	// 初始级别应保持 warn
	if got := GetLogLevel(); got != zap.WarnLevel {
		t.Fatalf("expect initial level=warn when file missing, got %s", got.String())
	}

	// 稍后创建文件
	time.Sleep(200 * time.Millisecond)
	if err := os.WriteFile(confPath, []byte("env_config:\n  log_level: debug\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// 等待回调生效（给足 5s，兼容 polling 周期 200ms 与 fsnotify 延迟）
	if !waitLogLevel(zap.DebugLevel, 5*time.Second) {
		t.Fatalf("expect level=debug after file created, got %s", GetLogLevel().String())
	}
}
