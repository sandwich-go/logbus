package logbus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sandwich-go/xconf-providers/xfile"
	"github.com/sandwich-go/xconf/kv"
	"go.uber.org/zap/zapcore"
)

// 环境变量 key，部署时默认注入的环境变量。
const (
	envKeyConfPathEnv = "sys_conf_path_env"
	envKeyCDService   = "sys_cd_service"
	// bizopsConfFile PMT 下发的业务运维配置文件名，放在 $sys_conf_path_env 目录下。
	bizopsConfFile = "ops_config.json"
)

// bizOpsConf 对应 ops_config.json 的结构。
// 字段级覆盖：service_config[<currentService>] 中存在的字段覆盖 env_config 中的对应字段。
//
//	{
//	  "env_config": { "log_level": "info" },
//	  "service_config": {
//	    "svcA": { "log_level": "debug" }
//	  }
//	}
//
// 使用标准库 encoding/json 做纯解析，避免引入 yaml/xconf 等可能做环境变量替换的逻辑。
type bizOpsConf struct {
	EnvConfig     logLevelConf            `json:"env_config"`
	ServiceConfig map[string]logLevelConf `json:"service_config"`
}

type logLevelConf struct {
	LogLevel string `json:"log_level"`
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

	confDir := strings.TrimSpace(os.Getenv(envKeyConfPathEnv))
	if confDir == "" {
		// 非 PMT 环境不启用，业务本地开发不受影响
		return
	}
	confPath := filepath.Join(confDir, bizopsConfFile)
	serviceName := strings.TrimSpace(os.Getenv(envKeyCDService))
	if dynamicLogLevelStarted {
		applyLogConfFromLoader(context.Background(), dynamicLogLevelLoader, confPath, serviceName)
		return
	}

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

	ctx := context.Background()
	applyLogConfFromLoader(ctx, loader, confPath, serviceName)

	loader.Watch(ctx, confPath, func(_ string, path string, content []byte) error {
		applyLogConfContent(path, content, serviceName)
		return nil
	})

	dynamicLogLevelLoader = loader
	dynamicLogLevelStarted = true

	InfoWithChannel(SERVERLOG, "logbus dynamic loglevel enabled",
		String("path", confPath), String("service", serviceName))
}

func applyLogConfFromLoader(ctx context.Context, loader kv.Loader, confPath, serviceName string) {
	if loader == nil {
		return
	}
	// 初始读取：PollingMode 下读取失败会吞掉 err（见 xfile.GetImplement），
	// 此时文件可能尚未下发，跳过即可；待文件出现后由 Watch 回调应用。
	if data, gErr := loader.Get(ctx, confPath); gErr != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: read initial conf failed",
			String("path", confPath), String("err", gErr.Error()))
	} else if len(data) > 0 {
		applyLogConfContent(confPath, data, serviceName)
	}
}

// closeDynamicLogLevel 关闭 loader，释放文件监听资源。
//
// 历史背景：xconf <= v0.3.28 的 kv.Common.Done 在 New 时未 make，
// Close 中 close(c.Done) 会 panic: close of nil channel。
// 该缺陷在 xconf v0.3.29 已修复（New 中 make Done + Close 用 sync.Once 幂等）。
//
// 防御性保留 recover：避免日后再有上游回归或第三方实现 kv.Loader 时
// Close 抛 panic 直接打挂调用方进程；正常路径不会被触发，仅在 debug 级日志记录。
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
		defer func() {
			if r := recover(); r != nil {
				DebugWithChannel(SERVERLOG, "logbus dynamic loglevel: loader.Close panicked (recovered, defensive)",
					String("reason", fmt.Sprintf("%v", r)))
			}
		}()
		_ = loader.Close(context.Background())
	}()
}

// applyLogConfContent 解析 json 内容并根据服务名决定最终 LogLevel。
// 使用标准库 encoding/json，不走 xconf（xconf 会做 ${ENV} 替换，可能污染配置）。
//
// 决策规则（按优先级从高到低，前一个解析失败/为空则尝试下一个）：
//  1. service_config[<serviceName>].log_level
//  2. env_config.log_level
//
// 全部为空 → 静默返回（保持当前 level 不动）。
// 候选项非法仅打 warning 不影响其它候选项继续尝试，避免 service 写错连累 env 的合法配置。
func applyLogConfContent(path string, content []byte, serviceName string) {
	if len(content) == 0 {
		return
	}
	var cfg bizOpsConf
	if err := json.Unmarshal(content, &cfg); err != nil {
		WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: unmarshal json failed",
			String("path", path), String("err", err.Error()))
		return
	}

	type candidate struct {
		from  string // 用于日志，标识来源
		value string
	}
	candidates := make([]candidate, 0, 2)
	if serviceName != "" {
		if sc, ok := cfg.ServiceConfig[serviceName]; ok && strings.TrimSpace(sc.LogLevel) != "" {
			candidates = append(candidates, candidate{from: "service_config." + serviceName, value: sc.LogLevel})
		}
	}
	if strings.TrimSpace(cfg.EnvConfig.LogLevel) != "" {
		candidates = append(candidates, candidate{from: "env_config", value: cfg.EnvConfig.LogLevel})
	}
	if len(candidates) == 0 {
		return
	}

	for _, c := range candidates {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(strings.ToLower(strings.TrimSpace(c.value)))); err != nil {
			WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: invalid log_level, try next candidate",
				String("path", path), String("from", c.from), String("level", c.value), String("err", err.Error()))
			continue
		}
		if cur := GetLogLevel(); cur == level {
			return
		}
		SetLogLevel(level)
		InfoWithChannel(SERVERLOG, "logbus dynamic loglevel updated",
			String("path", path), String("from", c.from), String("service", serviceName), String("level", level.String()))
		return
	}
	// 所有候选项均非法
	WarnWithChannel(SERVERLOG, "logbus dynamic loglevel: all candidates invalid, keep current level",
		String("path", path), String("service", serviceName), String("current_level", GetLogLevel().String()))
}
