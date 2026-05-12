package logbus

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sandwich-go/xconf-providers/xfile"
	"github.com/sandwich-go/xconf/kv"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// 环境变量 key，部署时默认注入的环境变量。
const (
	envKeyConfPathEnv = "sys_conf_path_env"
	envKeyCDService   = "sys_cd_service"
	// logConfFile 配置文件名，放在 sys_conf_path_env 目录下。
	logConfFile = "bizops.yaml"
)

// bizOpsConf 对应 bizops.yaml 的结构。
// 字段级覆盖：service_config[<currentService>] 中存在的字段覆盖 env_config 中的对应字段。
//
//	env_config:
//	  log_level: info
//	service_config:
//	  svcA:
//	    log_level: debug
type bizOpsConf struct {
	EnvConfig     logLevelConf            `yaml:"env_config"`
	ServiceConfig map[string]logLevelConf `yaml:"service_config"`
}

type logLevelConf struct {
	LogLevel string `yaml:"log_level"`
}

var (
	dynamicLogLevelMu      sync.Mutex
	dynamicLogLevelLoader  kv.Loader
	dynamicLogLevelStarted bool
)

// enableDynamicLogLevel 在 Init() 末尾调用：
//   - 若 sys_conf_path_env 环境变量不存在或目标文件无法定位，直接返回（不报错，仅 warning）
//   - 否则通过 xfile.Loader 读取初始配置并 Watch 文件变化，变化时按 “服务级字段覆盖环境级字段” 决策 LogLevel。
//
// 只初始化一次，重复调用安全。
func enableDynamicLogLevel() {
	dynamicLogLevelMu.Lock()
	defer dynamicLogLevelMu.Unlock()
	if dynamicLogLevelStarted {
		return
	}

	confDir := strings.TrimSpace(os.Getenv(envKeyConfPathEnv))
	if confDir == "" {
		// 非 PMT 环境不启用，业务本地开发不受影响
		return
	}
	confPath := filepath.Join(confDir, logConfFile)

	// 强制使用 PollingMode：
	//   1. 文件从无到有场景下 fsnotify 的 Add() 会直接失败，后续即便文件创建也不会触发事件；
	//      poller 则允许监听尚未存在的路径，文件出现时产生 Create 事件。
	//   2. K8s configmap 挂载通过 symlink 原子切换 ..data 目录，fsnotify 监听挂载点文件的
	//      写事件在部分 K8s 版本上并不可靠；polling 按 mtime+size 对比稳定可靠。
	loader, err := xfile.New(
		xfile.WithPollingMode(true),
		xfile.WithLogDebug(func(s string) { DebugWithChannel(SERVERLOG, "logbus dynamic loglevel: "+s) }),
		xfile.WithLogWarning(func(s string) { WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: "+s) }),
	)
	if err != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: new xfile loader failed",
			String("err", err.Error()))
		return
	}

	serviceName := strings.TrimSpace(os.Getenv(envKeyCDService))

	ctx := context.Background()

	// 初始读取：PollingMode 下读取失败会吞掉 err（见 xfile.GetImplement），
	// 此时文件可能尚未下发，跳过即可；待文件出现后由 Watch 回调应用。
	if data, gErr := loader.Get(ctx, confPath); gErr != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: read initial conf failed",
			String("path", confPath), String("err", gErr.Error()))
	} else if len(data) > 0 {
		applyLogConfContent(confPath, data, serviceName)
	}

	loader.Watch(ctx, confPath, func(_ string, path string, content []byte) error {
		applyLogConfContent(path, content, serviceName)
		return nil
	})

	dynamicLogLevelLoader = loader
	dynamicLogLevelStarted = true

	InfoWithChannel(SERVERLOG, "logbus dynamic loglevel enabled",
		String("path", confPath), String("service", serviceName))
}

// closeDynamicLogLevel 关闭 loader，释放文件监听资源。
// 使用 recover 兜底，避免上游 xconf/kv.Common.Close 在 Done 通道未初始化场景下 panic。
func closeDynamicLogLevel() {
	dynamicLogLevelMu.Lock()
	defer dynamicLogLevelMu.Unlock()
	if !dynamicLogLevelStarted || dynamicLogLevelLoader == nil {
		return
	}
	loader := dynamicLogLevelLoader
	dynamicLogLevelLoader = nil
	dynamicLogLevelStarted = false
	func() {
		defer func() { _ = recover() }()
		_ = loader.Close(context.Background())
	}()
}

// applyLogConfContent 解析 yaml 内容并根据服务名决定最终 LogLevel。
func applyLogConfContent(path string, content []byte, serviceName string) {
	if len(content) == 0 {
		return
	}
	var cfg bizOpsConf
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: unmarshal yaml failed",
			String("path", path), String("err", err.Error()))
		return
	}

	// 字段级覆盖：service_config[<serviceName>].log_level > env_config.log_level
	levelStr := cfg.EnvConfig.LogLevel
	if serviceName != "" {
		if sc, ok := cfg.ServiceConfig[serviceName]; ok && sc.LogLevel != "" {
			levelStr = sc.LogLevel
		}
	}
	levelStr = strings.TrimSpace(levelStr)
	if levelStr == "" {
		return
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(levelStr))); err != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: invalid log_level",
			String("path", path), String("level", levelStr), String("err", err.Error()))
		return
	}

	if cur := GetLogLevel(); cur == level {
		return
	}
	SetLogLevel(level)
	InfoWithChannel(SERVERLOG, "logbus dynamic loglevel updated",
		String("path", path), String("service", serviceName), String("level", level.String()))
}
