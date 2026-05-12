// Package main demonstrates logbus dynamic log level driven by PMT-managed configmap.
//
// 运行前确认当前目录下有 bizops.yaml，然后：
//
//	# 模拟 PMT 注入的环境变量
//	export sys_conf_path_env=$(pwd)
//	export sys_cd_service=service_a   # 当前服务名，与 bizops.yaml 内 service_config.<name> 对齐
//
//	go run .
//
// 程序启动后会每秒打印一条 Debug / Info / Warn / Error 日志，
// 同时观察日志输出变化，可以直接 `vim bizops.yaml` 修改 log_level，
// logbus 会在 <= 200ms 内应用新级别，无需重启进程。
package main

import (
	"os"
	"time"

	"github.com/sandwich-go/logbus"
	"go.uber.org/zap"
)

func main() {
	defer logbus.Close()

	// Init 末尾会自动调用 enableDynamicLogLevel()：
	//   - 若未设置 sys_conf_path_env 环境变量，本地运行不会启用 watch（零侵入）
	//   - 若已设置，则 watch $sys_conf_path_env/bizops.yaml，支持从无到有
	logbus.Init(logbus.NewConf(
		logbus.WithDev(true),
		logbus.WithDefaultChannel("DynamicLogLevel"),
		logbus.WithLogLevel(zap.WarnLevel), // 默认 warn，期望被 bizops.yaml 覆盖为 info/debug
	))

	if os.Getenv("sys_conf_path_env") == "" {
		logbus.Warn("sys_conf_path_env not set; dynamic log level is NOT enabled, set it and rerun to see effect")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			logbus.Debug("debug tick", logbus.Int("i", i), logbus.String("current_level", logbus.GetLogLevel().String()))
			logbus.Info("info tick", logbus.Int("i", i), logbus.String("current_level", logbus.GetLogLevel().String()))
			logbus.Warn("warn tick", logbus.Int("i", i), logbus.String("current_level", logbus.GetLogLevel().String()))
			if i > 0 && i%10 == 0 {
				logbus.Info("tip: try editing bizops.yaml `log_level` value, changes take effect within ~200ms")
			}
		}
	}
}
